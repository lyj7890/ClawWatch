package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := loadConfig()
	hub := NewHub()

	go hub.Run()

	mux := http.NewServeMux()

	// WebSocket 端点
	mux.HandleFunc("/ws/agent", handleAgentWS(hub, cfg))
	mux.HandleFunc("/ws/console", handleConsoleWS(hub, cfg))

	// REST API
	mux.HandleFunc("/api/agents", handleAgentsList(hub))
	mux.HandleFunc("/api/health", handleHealth())

	// 静态页面（简单 console UI）
	mux.HandleFunc("/", handleIndex(cfg))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("[hub] ClawWatch Hub starting on %s", addr)
	log.Printf("[hub] Agent WS:   ws://0.0.0.0%s/ws/agent?token=<token>&agentId=<id>", addr)
	log.Printf("[hub] Console WS: ws://0.0.0.0%s/ws/console?token=<token>&subscribe=*", addr)

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

func handleAgentsList(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub.mu.RLock()
		agents := make([]map[string]interface{}, 0, len(hub.agents))
		for id, a := range hub.agents {
			agents = append(agents, map[string]interface{}{
				"id":       id,
				"lastSeen": a.LastSeen.UnixMilli(),
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
			"version": "0.1.0",
		})
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
