package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

const (
	protocolVersion   = 3
	agentVersion      = "0.3.1"
	pollInterval      = time.Second
	reconnectInterval = 3 * time.Second
	maxPendingEvents  = 10000
)

type config struct {
	HubURL           string
	WatchDir         string
	TrajectoryRedact bool
}

type hostInfo struct {
	ID       string   `json:"id"`
	Hostname string   `json:"hostname"`
	IPs      []string `json:"ips"`
	OS       string   `json:"os"`
	Arch     string   `json:"arch"`
}

type openClawAgentInfo struct {
	ID           string `json:"id"`
	SessionCount int    `json:"sessionCount"`
	LastActivity int64  `json:"lastActivity,omitempty"`
}

type sessionInfo struct {
	AgentID   string `json:"agentId"`
	SessionID string `json:"sessionId"`
	Path      string `json:"path"`
	MTime     int64  `json:"mtime"`
	Size      int64  `json:"size"`
}

type fileState struct {
	Offset    int64
	Remainder []byte
	AgentID   string
	SessionID string
}

type agent struct {
	cfg      config
	host     hostInfo
	mu       sync.Mutex
	conn     *websocket.Conn
	pending  [][]byte
	files    map[string]*fileState
	stopping chan struct{}
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	hostID := loadOrCreateHostID(homeDir(), hostname())
	a := &agent{
		cfg:      cfg,
		host:     hostInfo{ID: hostID, Hostname: displayHostname(hostname(), hostID), IPs: localIPs(), OS: runtime.GOOS, Arch: runtime.GOARCH},
		files:    make(map[string]*fileState),
		stopping: make(chan struct{}),
	}

	log.Printf("[agent] host=%s hostname=%s ips=%s watch=%s", a.host.ID, a.host.Hostname, strings.Join(a.host.IPs, ","), cfg.WatchDir)
	go a.connectLoop()
	go a.scanLoop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	close(a.stopping)
	a.closeConnection()
}

func loadConfig() (config, error) {
	cfg := config{}
	flag.StringVar(&cfg.HubURL, "hub", "", "Hub base WebSocket URL, for example wss://clawhub.example.com")
	flag.StringVar(&cfg.WatchDir, "dir", filepath.Join(homeDir(), ".openclaw", "agents"), "OpenClaw agents directory")
	flag.BoolVar(&cfg.TrajectoryRedact, "trajectory-redact", true, "Redact sensitive trajectory data before reporting")
	flag.Parse()
	if cfg.HubURL == "" {
		return cfg, fmt.Errorf("--hub is required")
	}
	return cfg, nil
}

func hostname() string {
	value, err := os.Hostname()
	if err != nil || value == "" {
		return "unknown-host"
	}
	return value
}

func displayHostname(value, hostID string) string {
	switch strings.ToLower(value) {
	case "", "localhost", "localhost.localdomain", "unknown-host":
		return hostID
	default:
		return value
	}
}

func homeDir() string {
	value, _ := os.UserHomeDir()
	return value
}

func loadOrCreateHostID(home, hostname string) string {
	path := filepath.Join(home, ".clawwatch", "host-id")
	if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) != "" {
		return strings.TrimSpace(string(data))
	}
	random := make([]byte, 6)
	if _, err := rand.Read(random); err != nil {
		return hostname
	}
	id := fmt.Sprintf("%s-%x", hostname, random)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err == nil {
		_ = os.WriteFile(path, []byte(id+"\n"), 0o600)
	}
	return id
}

func localIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var ips []string
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}
		ips = append(ips, ip.String())
	}
	sort.Strings(ips)
	return ips
}

func (a *agent) connectLoop() {
	for {
		select {
		case <-a.stopping:
			return
		default:
		}

		u, err := url.Parse(strings.TrimRight(a.cfg.HubURL, "/") + "/ws/agent")
		if err != nil {
			log.Printf("[agent] invalid Hub URL: %v", err)
			return
		}
		query := u.Query()
		query.Set("agentId", a.host.ID)
		u.RawQuery = query.Encode()

		dialer := *websocket.DefaultDialer
		dialer.Proxy = nil
		conn, _, err := dialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("[agent] connect failed: %v", err)
			if !a.wait(reconnectInterval) {
				return
			}
			continue
		}

		a.mu.Lock()
		a.conn = conn
		a.mu.Unlock()
		log.Printf("[agent] connected to %s", a.cfg.HubURL)
		a.sendHelloAndInventory()
		a.flushPending()

		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if msgType == websocket.TextMessage {
				a.handleMessage(data)
			}
		}
		a.mu.Lock()
		if a.conn == conn {
			a.conn = nil
		}
		a.mu.Unlock()
		_ = conn.Close()
		log.Printf("[agent] disconnected; reconnecting")
	}
}

func (a *agent) wait(d time.Duration) bool {
	select {
	case <-a.stopping:
		return false
	case <-time.After(d):
		return true
	}
}

func (a *agent) closeConnection() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn != nil {
		_ = a.conn.Close()
		a.conn = nil
	}
}

func (a *agent) handleMessage(data []byte) {
	var msg map[string]any
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[agent] failed to parse message: %v", err)
		return
	}
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}
	switch msgType {
	case "fetch_session":
		a.handleFetchSession(msg)
	case "agent_hello_ack":
		if token, ok := msg["consoleToken"].(string); ok && token != "" {
			log.Printf("[agent] Console token for this host: %s", token)
		}
	default:
		// Ignore unknown message types
	}
}

func (a *agent) handleFetchSession(msg map[string]any) {
	requestID, _ := msg["requestId"].(string)
	openclawAgentID, _ := msg["openclawAgentId"].(string)
	sessionID, _ := msg["sessionId"].(string)

	// Default tail 100 lines, allow override
	tailLines := 100
	if v, ok := msg["tailLines"].(float64); ok && v > 0 && v <= 500 {
		tailLines = int(v)
	}

	if requestID == "" || openclawAgentID == "" || sessionID == "" {
		log.Printf("[agent] fetch_session missing required fields")
		return
	}

	// Construct session file path
	sessionPath := filepath.Join(a.cfg.WatchDir, openclawAgentID, "sessions", sessionID+".jsonl")

	// Check file exists
	_, err := os.Stat(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			a.send(map[string]any{
				"type":      "session_history_error",
				"requestId": requestID,
				"error":     "session not found",
			})
			return
		}
		a.send(map[string]any{
			"type":      "session_history_error",
			"requestId": requestID,
			"error":     fmt.Sprintf("stat failed: %v", err),
		})
		return
	}

	// Read only the last N lines (tail) to avoid memory exhaustion
	lines, err := tailFile(sessionPath, tailLines)
	if err != nil {
		log.Printf("[agent] read session file failed: %v", err)
		a.send(map[string]any{
			"type":      "session_history_error",
			"requestId": requestID,
			"error":     fmt.Sprintf("read failed: %v", err),
		})
		return
	}

	// Send successful response
	a.send(map[string]any{
		"type":            "session_history",
		"requestId":       requestID,
		"openclawAgentId": openclawAgentID,
		"sessionId":       sessionID,
		"lines":           lines,
		"timestamp":       time.Now().UnixMilli(),
	})
}

// tailFile reads the last N valid JSON lines from a jsonl file
func tailFile(path string, n int) ([]json.RawMessage, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read all lines but keep only last N (streaming to avoid loading huge files fully)
	// For files under 2MB, read fully; for larger files, seek to approximate tail
	info, _ := f.Stat()
	var scanner *bufio.Scanner
	if info.Size() > 2*1024*1024 {
		// Seek to last 2MB for tail approximation
		offset := info.Size() - 2*1024*1024
		f.Seek(offset, 0)
		// Discard partial first line
		br := bufio.NewReader(f)
		br.ReadBytes('\n')
		scanner = bufio.NewScanner(br)
	} else {
		scanner = bufio.NewScanner(f)
	}
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

	// Ring buffer for last N lines
	ring := make([]json.RawMessage, 0, n)
	for scanner.Scan() {
		raw := append([]byte(nil), scanner.Bytes()...)
		if len(raw) > 0 && json.Valid(raw) {
			if len(ring) >= n {
				ring = ring[1:]
			}
			ring = append(ring, raw)
		}
	}
	return ring, scanner.Err()
}

func (a *agent) scanLoop() {
	a.scan(true)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-a.stopping:
			return
		case <-ticker.C:
			a.scan(false)
		}
	}
}

func (a *agent) scan(initial bool) {
	sessions, err := discoverSessions(a.cfg.WatchDir)
	if err != nil {
		log.Printf("[agent] scan failed: %v", err)
		return
	}
	for _, session := range sessions {
		state, exists := a.files[session.Path]
		if !exists {
			offset := int64(0)
			if initial {
				offset = session.Size
			}
			state = &fileState{Offset: offset, AgentID: session.AgentID, SessionID: session.SessionID}
			a.files[session.Path] = state
			if !initial {
				a.send(map[string]any{
					"type": "new_session", "protocolVersion": protocolVersion,
					"host": a.host, "openclawAgentId": session.AgentID, "sessionId": session.SessionID,
					"timestamp": time.Now().UnixMilli(),
				})
			}
		}
		if session.Size > state.Offset {
			a.readNewLines(session.Path, state)
		}
	}
	trajectories, err := discoverTrajectories(a.cfg.WatchDir)
	if err != nil {
		log.Printf("[agent] trajectory scan failed: %v", err)
		return
	}
	for _, trajectory := range trajectories {
		state, exists := a.files[trajectory.Path]
		if !exists {
			offset := int64(0)
			if initial {
				offset = trajectory.Size
			}
			state = &fileState{Offset: offset, AgentID: trajectory.AgentID, SessionID: trajectory.SessionID}
			a.files[trajectory.Path] = state
		}
		if trajectory.Size > state.Offset {
			a.readNewTrajectoryLines(trajectory.Path, state)
		}
	}
}

func discoverSessions(root string) ([]sessionInfo, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var sessions []sessionInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		agentID := entry.Name()
		sessionDir := filepath.Join(root, agentID, "sessions")
		files, err := os.ReadDir(sessionDir)
		if err != nil {
			continue
		}
		for _, file := range files {
			name := file.Name()
			if file.IsDir() || !strings.HasSuffix(name, ".jsonl") || strings.Contains(name, ".trajectory") || strings.Contains(name, ".checkpoint") {
				continue
			}
			info, err := file.Info()
			if err != nil {
				continue
			}
			sessions = append(sessions, sessionInfo{
				AgentID: agentID, SessionID: strings.TrimSuffix(name, ".jsonl"),
				Path: filepath.Join(sessionDir, name), MTime: info.ModTime().UnixMilli(), Size: info.Size(),
			})
		}
	}
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].MTime > sessions[j].MTime })
	// Only track the most recent 20 sessions globally to limit resource usage
	const maxSessions = 20
	if len(sessions) > maxSessions {
		sessions = sessions[:maxSessions]
	}
	return sessions, nil
}

func summarizeAgents(root string, sessions []sessionInfo) []openClawAgentInfo {
	// First, discover all agent directories
	byID := make(map[string]*openClawAgentInfo)
	entries, _ := os.ReadDir(root)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Only include if it has a sessions/ subdirectory
		sessionDir := filepath.Join(root, entry.Name(), "sessions")
		if _, err := os.Stat(sessionDir); err == nil {
			byID[entry.Name()] = &openClawAgentInfo{ID: entry.Name()}
		}
	}
	// Then enrich with session counts from tracked sessions
	for _, session := range sessions {
		info := byID[session.AgentID]
		if info == nil {
			info = &openClawAgentInfo{ID: session.AgentID}
			byID[session.AgentID] = info
		}
		info.SessionCount++
		if session.MTime > info.LastActivity {
			info.LastActivity = session.MTime
		}
	}
	agents := make([]openClawAgentInfo, 0, len(byID))
	for _, info := range byID {
		agents = append(agents, *info)
	}
	sort.Slice(agents, func(i, j int) bool { return agents[i].ID < agents[j].ID })
	return agents
}

func (a *agent) sendHelloAndInventory() {
	sessions, _ := discoverSessions(a.cfg.WatchDir)
	a.send(map[string]any{
		"type": "agent_hello", "protocolVersion": protocolVersion, "agentVersion": agentVersion,
		"host": a.host, "agents": summarizeAgents(a.cfg.WatchDir, sessions), "timestamp": time.Now().UnixMilli(),
	})
	a.send(map[string]any{
		"type": "session_list", "protocolVersion": protocolVersion, "host": a.host,
		"sessions": sessions, "timestamp": time.Now().UnixMilli(),
	})
}

func (a *agent) readNewLines(path string, state *fileState) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	if _, err := file.Seek(state.Offset, io.SeekStart); err != nil {
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	state.Offset += int64(len(data))
	data = append(state.Remainder, data...)

	var lines []json.RawMessage
	lastNewline := strings.LastIndexByte(string(data), '\n')
	if lastNewline < 0 {
		state.Remainder = append(state.Remainder[:0], data...)
		return
	}
	complete := data[:lastNewline]
	state.Remainder = append(state.Remainder[:0], data[lastNewline+1:]...)
	scanner := bufio.NewScanner(strings.NewReader(string(complete)))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		raw := append([]byte(nil), scanner.Bytes()...)
		if json.Valid(raw) {
			lines = append(lines, raw)
		}
	}
	if len(lines) == 0 {
		return
	}
	a.send(map[string]any{
		"type": "log", "protocolVersion": protocolVersion, "host": a.host,
		"openclawAgentId": state.AgentID, "sessionId": state.SessionID,
		"lines": lines, "timestamp": time.Now().UnixMilli(),
	})
}

func (a *agent) send(message any) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn != nil {
		_ = a.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := a.conn.WriteMessage(websocket.TextMessage, data); err == nil {
			return
		}
		_ = a.conn.Close()
		a.conn = nil
	}
	a.enqueueLocked(data)
}

func (a *agent) enqueueLocked(data []byte) {
	if len(a.pending) >= maxPendingEvents {
		copy(a.pending, a.pending[1:])
		a.pending = a.pending[:maxPendingEvents-1]
	}
	a.pending = append(a.pending, data)
}

func (a *agent) flushPending() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn == nil {
		return
	}
	for len(a.pending) > 0 {
		if err := a.conn.WriteMessage(websocket.TextMessage, a.pending[0]); err != nil {
			_ = a.conn.Close()
			a.conn = nil
			return
		}
		a.pending = a.pending[1:]
	}
}
