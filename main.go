package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	cfg := loadConfig()
	store, err := NewMessageStore(cfg.StoragePath, cfg.HistoryLimit)
	if err != nil {
		log.Fatalf("[hub] initialize message store: %v", err)
	}
	metrics, err := NewMessageStore(cfg.MetricsPath, cfg.MetricsLimit, "trajectory")
	if err != nil {
		log.Fatalf("[hub] initialize metrics store: %v", err)
	}
	tokens, err := NewTokenStore(cfg.TokenStorePath)
	if err != nil {
		log.Fatalf("[hub] initialize token store: %v", err)
	}
	hub := NewHub(store, metrics, tokens)

	go hub.Run()

	mux := http.NewServeMux()

	// WebSocket 端点
	mux.HandleFunc("/ws/agent", handleAgentWS(hub, cfg))
	mux.HandleFunc("/ws/console", handleConsoleWS(hub, cfg))

	// REST API
	mux.HandleFunc("/api/agents", handleAgentsList(hub, cfg))
	mux.HandleFunc("/api/history", handleHistory(hub, cfg, store))
	mux.HandleFunc("/api/trajectory", handleTrajectoryHistory(hub, cfg, metrics))
	mux.HandleFunc("/api/session-history", handleSessionHistory(hub, cfg))
	mux.HandleFunc("/api/health", handleHealth())

	// 静态页面（简单 console UI）
	mux.HandleFunc("/", handleIndex(cfg))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("[hub] ClawWatch Hub starting on %s", addr)
	log.Printf("[hub] Agent WS:   ws://0.0.0.0%s/ws/agent?token=<token>&agentId=<id>", addr)
	log.Printf("[hub] Console WS: ws://0.0.0.0%s/ws/console?token=<token>&subscribe=*", addr)
	log.Printf("[hub] History:    %s (limit=%d)", cfg.StoragePath, cfg.HistoryLimit)
	log.Printf("[hub] Trajectory: %s (limit=%d)", cfg.MetricsPath, cfg.MetricsLimit)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // WebSocket 需要长连接
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[hub] server error: %v", err)
	}
}

// consoleScope 表示一个 REST 请求凭 token 解析出的可见范围。
type consoleScope struct {
	admin          bool
	allowedAgentID string
	authorized     bool
}

func (s consoleScope) canSee(agentID string) bool {
	if s.admin {
		return true
	}
	return s.allowedAgentID != "" && s.allowedAgentID == agentID
}

// resolveScope 从请求的 token 查询縁7的可见范围，与 console WS 的逻辑保持一致。
func resolveScope(hub *Hub, cfg *Config, r *http.Request) consoleScope {
	token := r.URL.Query().Get("token")
	scope := consoleScope{}

	switch {
	case cfg.ConsoleToken != "" && token == cfg.ConsoleToken:
		scope.admin = true
	case hub.tokens != nil:
		if agentID, ok := hub.tokens.AgentIDForToken(token); ok {
			scope.allowedAgentID = agentID
		}
	}

	if cfg.ConsoleToken != "" {
		scope.authorized = scope.admin || scope.allowedAgentID != ""
	} else {
		if scope.allowedAgentID == "" {
			scope.admin = true
		}
		scope.authorized = true
	}
	return scope
}

// denyAgentID 是一个不可能匹配任何真实 agentId 的哨兵值，表示请求越权。
const denyAgentID = "\x00__deny__"

// scopedAgentID 根据可见范围修正请求的 agentId 参数。
// admin：原样返回（空/"*" = 全部）。
// 非 admin：忽略请求参数，强制锁定到 allowedAgentID；若显式请求了其他 agentId 则拒绝。
func scopedAgentID(scope consoleScope, requested string) string {
	if scope.admin {
		return requested
	}
	if requested != "" && requested != "*" && requested != scope.allowedAgentID {
		return denyAgentID
	}
	return scope.allowedAgentID
}

func handleTrajectoryHistory(hub *Hub, cfg *Config, store *MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scope := resolveScope(hub, cfg, r)
		if !scope.authorized {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		limit := 500
		if value := r.URL.Query().Get("limit"); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		agentID := scopedAgentID(scope, r.URL.Query().Get("agentId"))
		if agentID == denyAgentID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		messages := store.List(agentID, limit)
		for index, raw := range messages {
			var event map[string]interface{}
			if json.Unmarshal(raw, &event) == nil {
				delete(event, "metrics")
				if cleaned, err := json.Marshal(event); err == nil {
					messages[index] = cleaned
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"messages": messages,
			"count":    len(messages),
		})
	}
}

func handleHistory(hub *Hub, cfg *Config, store *MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scope := resolveScope(hub, cfg, r)
		if !scope.authorized {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		limit := 500
		if value := r.URL.Query().Get("limit"); value != "" {
			if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		agentID := scopedAgentID(scope, r.URL.Query().Get("agentId"))
		if agentID == denyAgentID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		messages := store.List(agentID, limit)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"messages": messages,
			"count":    len(messages),
		})
	}
}

func handleAgentsList(hub *Hub, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scope := resolveScope(hub, cfg, r)
		if !scope.authorized {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		hub.mu.RLock()
		agents := make([]map[string]interface{}, 0, len(hub.agents))
		for id, a := range hub.agents {
			if !scope.canSee(id) {
				continue
			}
			agents = append(agents, map[string]interface{}{
				"id":              id,
				"hostname":        a.Hostname,
				"hostIPs":         a.HostIPs,
				"os":              a.OS,
				"arch":            a.Arch,
				"agentVersion":    a.AgentVersion,
				"protocolVersion": a.ProtocolVersion,
				"openclawAgents":  a.OpenClawAgents,
				"sessions":        a.Sessions,
				"lastSeen":        a.LastSeen.UnixMilli(),
			})
		}
		hub.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"agents": agents,
			"count":  len(agents),
		})
	}
}

func handleHealth() http.HandlerFunc {
	start := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"uptime":  time.Since(start).String(),
			"version": "0.2.0",
		})
	}
}

func handleSessionHistory(hub *Hub, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scope := resolveScope(hub, cfg, r)
		if !scope.authorized {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// 获取参数
		agentID := r.URL.Query().Get("agentId")
		openclawAgentID := r.URL.Query().Get("openclawAgentId")
		sessionID := r.URL.Query().Get("sessionId")

		if agentID == "" || openclawAgentID == "" || sessionID == "" {
			http.Error(w, "agentId, openclawAgentId, and sessionId are required", http.StatusBadRequest)
			return
		}

		// 权限校验：非 admin 只能访问自己 token 对应的 agent
		if !scope.canSee(agentID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// 查找 agent
		hub.mu.RLock()
		agent, exists := hub.agents[agentID]
		hub.mu.RUnlock()

		if !exists {
			http.Error(w, fmt.Sprintf("agent %s not found", agentID), http.StatusNotFound)
			return
		}

		// 生成 requestId
		requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

		// 创建响应 channel
		respChan := make(chan []byte, 1)
		hub.pendingMu.Lock()
		hub.pendingRequests[requestID] = respChan
		hub.pendingMu.Unlock()

		// 构造请求消息
		req := map[string]interface{}{
			"type":            "fetch_session",
			"requestId":       requestID,
			"openclawAgentId": openclawAgentID,
			"sessionId":       sessionID,
		}
		reqData, err := json.Marshal(req)
		if err != nil {
			http.Error(w, "failed to marshal request", http.StatusInternalServerError)
			return
		}

		// 发送给 agent
		select {
		case agent.Send <- reqData:
		case <-time.After(2 * time.Second):
			hub.pendingMu.Lock()
			delete(hub.pendingRequests, requestID)
			hub.pendingMu.Unlock()
			http.Error(w, "failed to send request to agent (timeout)", http.StatusGatewayTimeout)
			return
		}

		// 等待响应
		select {
		case resp := <-respChan:
			// 解析响应
			var respMsg struct {
				Type  string          `json:"type"`
				Data  json.RawMessage `json:"data"`
				Error string          `json:"error"`
			}
			if err := json.Unmarshal(resp, &respMsg); err != nil {
				http.Error(w, "failed to parse response", http.StatusInternalServerError)
				return
			}

			if respMsg.Type == "session_history_error" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": respMsg.Error,
				})
				return
			}

			// 返回成功响应
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(respMsg.Data)

		case <-time.After(10 * time.Second):
			hub.pendingMu.Lock()
			delete(hub.pendingRequests, requestID)
			hub.pendingMu.Unlock()
			http.Error(w, "timeout waiting for agent response", http.StatusGatewayTimeout)
		}
	}
}

func handleIndex(cfg *Config) http.HandlerFunc {
	html := buildIndexHTML(cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	}
}
