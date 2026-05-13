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
		// token 校验
		if cfg.ConsoleToken != "" && r.URL.Query().Get("token") != cfg.ConsoleToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		subscribe := r.URL.Query().Get("subscribe")
		if subscribe == "" {
			subscribe = "*"
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[console] upgrade error: %v", err)
			return
		}

		console := &ConsoleConn{
			ID:        generateID("console"),
			Subscribe: subscribe,
			Send:      make(chan []byte, 512),
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
