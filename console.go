package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// handleConsoleWS 处理 console 的 WebSocket 连接
// ws://hub:4848/ws/console?token=xxx&subscribe=my-mac  (subscribe="*" 订阅全部)
func handleConsoleWS(hub *Hub, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		// 解析 token 以确定该 console 的可见范围
		admin := false
		allowedAgentID := ""
		switch {
		case cfg.ConsoleToken != "" && token == cfg.ConsoleToken:
			// 全局 ConsoleToken = admin，可见全部 agent（向后兼容 + 管理用途）
			admin = true
		case hub.tokens != nil:
			if agentID, ok := hub.tokens.AgentIDForToken(token); ok {
				allowedAgentID = agentID
			}
		}

		// 权限校验：
		// - 若配置了全局 ConsoleToken，则必须是 admin 或匹配某个 per-agent token
		// - 若未配置全局 ConsoleToken，从宝未带 token 也允许（开放模式 = admin）
		if cfg.ConsoleToken != "" {
			if !admin && allowedAgentID == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		} else {
			if allowedAgentID == "" {
				// 未配置全局 token 且 token 不匹配任何 agent：视为 admin（开放模式）
				admin = true
			}
		}

		subscribe := r.URL.Query().Get("subscribe")
		if subscribe == "" {
			subscribe = "*"
		}
		// 非 admin 的 console 强制锁定到其 token 对应的 agentId
		if !admin {
			subscribe = allowedAgentID
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[console] upgrade error: %v", err)
			return
		}

		console := &ConsoleConn{
			ID:             generateID("console"),
			Subscribe:      subscribe,
			Send:           make(chan []byte, 512),
			AllowedAgentID: allowedAgentID,
			Admin:          admin,
		}
		hub.registerConsole <- console

		// write pump (blocks via goroutine, read pump blocks main)
		go consoleWritePump(conn, console)
		consoleReadPump(conn, hub, console)

		hub.unregisterConsole <- console
	}
}

func consoleReadPump(conn *websocket.Conn, hub *Hub, console *ConsoleConn) {
	defer conn.Close()
	conn.SetReadLimit(4096)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[console] %s read error: %v", console.ID, err)
			}
			return
		}
	}
}

func consoleWritePump(conn *websocket.Conn, console *ConsoleConn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-console.Send:
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
