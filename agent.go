package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 180 * time.Second
	pingPeriod = 45 * time.Second
	maxMsgSize = 4 << 20 // 4MB, supports explicitly unredacted trajectory events
)

// handleAgentWS 处理 agent 的 WebSocket 连接
// ws://hub:4848/ws/agent?token=xxx&agentId=my-mac
func handleAgentWS(hub *Hub, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// token 校验
		if cfg.AgentToken != "" && r.URL.Query().Get("token") != cfg.AgentToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		agentID := r.URL.Query().Get("agentId")
		if agentID == "" {
			http.Error(w, "agentId required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[agent] upgrade error: %v", err)
			return
		}

		agent := &AgentConn{
			ID:       agentID,
			Send:     make(chan []byte, 256),
			LastSeen: time.Now(),
		}
		hub.registerAgent <- agent

		// 为该 agentId 生成（或复用）专属 console token，并通过 ack 告知 agent
		if hub.tokens != nil {
			if token, err := hub.tokens.GetOrCreate(agentID); err != nil {
				log.Printf("[agent] %s: failed to provision console token: %v", agentID, err)
			} else {
				ack := map[string]interface{}{
					"type":         "agent_hello_ack",
					"agentId":      agentID,
					"consoleToken": token,
					"timestamp":    time.Now().UnixMilli(),
				}
				if data, err := json.Marshal(ack); err == nil {
					select {
					case agent.Send <- data:
					default:
						log.Printf("[agent] %s: failed to enqueue agent_hello_ack", agentID)
					}
				}
			}
		}

		// write pump
		go agentWritePump(conn, agent)

		// read pump (blocks)
		agentReadPump(conn, hub, agent)

		hub.unregisterAgent <- agent
	}
}

func agentReadPump(conn *websocket.Conn, hub *Hub, agent *AgentConn) {
	defer conn.Close()
	conn.SetReadLimit(maxMsgSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		hub.mu.Lock()
		agent.LastSeen = time.Now()
		hub.mu.Unlock()
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[agent] %s read error: %v", agent.ID, err)
			}
			return
		}
		hub.mu.Lock()
		agent.LastSeen = time.Now()
		updateAgentMetadata(raw, agent)
		hub.mu.Unlock()

		// 检查是否是 session_history 或 session_history_error 响应
		if handleSessionHistoryResponse(hub, raw) {
			// 内部通信，不广播给 console
			continue
		}

		hub.mu.Lock()
		enriched := injectAgentMetadata(raw, agent)
		hub.mu.Unlock()
		// 给 raw 消息注入 agentId 字段后广播给 consoles
		hub.broadcast <- &broadcastMsg{agentID: agent.ID, data: enriched}
	}
}

func agentWritePump(conn *websocket.Conn, agent *AgentConn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-agent.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func updateAgentMetadata(raw []byte, agent *AgentConn) {
	var msg struct {
		Type            string `json:"type"`
		ProtocolVersion int    `json:"protocolVersion"`
		AgentVersion    string `json:"agentVersion"`
		Host            struct {
			Hostname string   `json:"hostname"`
			IPs      []string `json:"ips"`
			OS       string   `json:"os"`
			Arch     string   `json:"arch"`
		} `json:"host"`
		Agents   []OpenClawAgentInfo `json:"agents"`
		Sessions []SessionInfo       `json:"sessions"`
	}
	if json.Unmarshal(raw, &msg) != nil {
		return
	}

	if msg.Type == "agent_hello" {
		agent.Hostname = msg.Host.Hostname
		agent.HostIPs = msg.Host.IPs
		agent.OS = msg.Host.OS
		agent.Arch = msg.Host.Arch
		agent.AgentVersion = msg.AgentVersion
		agent.ProtocolVersion = msg.ProtocolVersion
		agent.OpenClawAgents = msg.Agents
	} else if msg.Type == "session_list" {
		// 更新 sessions 列表
		agent.Sessions = msg.Sessions
	}
}

// injectAgentMetadata adds authoritative connection identity while preserving old clients.
func injectAgentMetadata(raw []byte, agent *AgentConn) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		// 非 JSON，包装一层
		return buildRawWrapper(agent.ID, raw)
	}
	m["agentId"] = agent.ID
	if agent.Hostname != "" {
		m["hostname"] = agent.Hostname
		m["hostIPs"] = agent.HostIPs
	}
	if _, ok := m["timestamp"]; !ok {
		m["timestamp"] = time.Now().UnixMilli()
	}
	out, err := json.Marshal(m)
	if err != nil {
		return buildRawWrapper(agent.ID, raw)
	}
	return out
}

func buildRawWrapper(agentID string, raw []byte) []byte {
	m := map[string]interface{}{
		"type":      "raw",
		"agentId":   agentID,
		"data":      string(raw),
		"timestamp": time.Now().UnixMilli(),
	}
	out, _ := json.Marshal(m)
	return out
}

// handleSessionHistoryResponse 处理 session_history 和 session_history_error 响应
// 返回 true 表示该消息已处理，不应广播给 console
func handleSessionHistoryResponse(hub *Hub, raw []byte) bool {
	var msg struct {
		Type      string          `json:"type"`
		RequestID string          `json:"requestId"`
		Data      json.RawMessage `json:"data"`
		Error     string          `json:"error"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return false
	}

	if msg.Type != "session_history" && msg.Type != "session_history_error" {
		return false
	}

	if msg.RequestID == "" {
		log.Printf("[agent] received %s without requestId", msg.Type)
		return true
	}

	hub.pendingMu.Lock()
	respChan, exists := hub.pendingRequests[msg.RequestID]
	if exists {
		delete(hub.pendingRequests, msg.RequestID)
	}
	hub.pendingMu.Unlock()

	if !exists {
		log.Printf("[agent] received %s for unknown requestId: %s", msg.Type, msg.RequestID)
		return true
	}

	// 发送响应到等待的 HTTP 请求
	select {
	case respChan <- raw:
	default:
		log.Printf("[agent] failed to send %s response for requestId: %s", msg.Type, msg.RequestID)
	}

	return true
}
