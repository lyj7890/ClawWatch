package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
	maxMsgSize = 1 << 20 // 1MB
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
		agent.LastSeen = time.Now()
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
		agent.LastSeen = time.Now()

		// 给 raw 消息注入 agentId 字段后广播给 consoles
		enriched := injectAgentID(raw, agent.ID)
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

// injectAgentID 向 JSON 消息中注入 agentId 字段
func injectAgentID(raw []byte, agentID string) []byte {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		// 非 JSON，包装一层
		return buildRawWrapper(agentID, raw)
	}
	m["agentId"] = agentID
	if _, ok := m["timestamp"]; !ok {
		m["timestamp"] = time.Now().UnixMilli()
	}
	out, err := json.Marshal(m)
	if err != nil {
		return buildRawWrapper(agentID, raw)
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
