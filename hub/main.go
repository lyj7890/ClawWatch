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

func buildIndexHTML(cfg *Config) string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>ClawWatch Hub</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: 'SF Mono', 'Fira Code', monospace; background: #0d1117; color: #c9d1d9; min-height: 100vh; }
header { background: #161b22; border-bottom: 1px solid #30363d; padding: 12px 20px; display: flex; align-items: center; gap: 12px; }
header h1 { font-size: 16px; color: #58a6ff; }
#status { font-size: 12px; color: #8b949e; }
#status.connected { color: #3fb950; }
#status.disconnected { color: #f85149; }
.container { display: grid; grid-template-columns: 220px 1fr; height: calc(100vh - 49px); }
.sidebar { background: #161b22; border-right: 1px solid #30363d; overflow-y: auto; }
.sidebar-title { padding: 12px 16px; font-size: 11px; color: #8b949e; text-transform: uppercase; letter-spacing: 1px; border-bottom: 1px solid #30363d; }
.agent-item { padding: 10px 16px; cursor: pointer; border-bottom: 1px solid #21262d; display: flex; align-items: center; gap: 8px; font-size: 13px; }
.agent-item:hover { background: #21262d; }
.agent-item.active { background: #1f2937; border-left: 2px solid #58a6ff; }
.dot { width: 8px; height: 8px; border-radius: 50%; background: #8b949e; flex-shrink: 0; }
.dot.online { background: #3fb950; }
.main { display: flex; flex-direction: column; }
#connect-bar { padding: 10px 16px; background: #161b22; border-bottom: 1px solid #30363d; display: flex; gap: 8px; align-items: center; }
#connect-bar input { flex: 1; background: #0d1117; border: 1px solid #30363d; color: #c9d1d9; padding: 6px 10px; border-radius: 6px; font-size: 13px; font-family: inherit; }
#connect-bar button { background: #238636; color: #fff; border: none; padding: 6px 14px; border-radius: 6px; cursor: pointer; font-size: 13px; }
#connect-bar button:hover { background: #2ea043; }
#log { flex: 1; overflow-y: auto; padding: 12px 16px; font-size: 12px; line-height: 1.6; }
.log-entry { margin-bottom: 4px; padding: 4px 8px; border-radius: 4px; word-break: break-all; }
.log-entry.log { color: #c9d1d9; }
.log-entry.agent_status { color: #79c0ff; background: #0d2137; }
.log-entry.error { color: #f85149; background: #2a0f0f; }
.log-entry.system { color: #8b949e; font-style: italic; }
.ts { color: #484f58; margin-right: 8px; }
</style>
</head>
<body>
<header>
  <h1>🦀 ClawWatch Hub</h1>
  <span id="status" class="disconnected">● Disconnected</span>
</header>
<div class="container">
  <div class="sidebar">
    <div class="sidebar-title">Agents</div>
    <div id="agent-list"></div>
  </div>
  <div class="main">
    <div id="connect-bar">
      <input id="token-input" type="password" placeholder="Console Token (留空跳过)">
      <input id="subscribe-input" type="text" placeholder="Subscribe agentId (* = 全部)" value="*">
      <button onclick="connect()">Connect</button>
    </div>
    <div id="log"></div>
  </div>
</div>
<script>
let ws = null;
const agents = {};

function ts() {
  return new Date().toLocaleTimeString('zh-CN', {hour12: false});
}

function appendLog(text, cls='log') {
  const log = document.getElementById('log');
  const div = document.createElement('div');
  div.className = 'log-entry ' + cls;
  div.innerHTML = '<span class="ts">' + ts() + '</span>' + escHtml(text);
  log.appendChild(div);
  log.scrollTop = log.scrollHeight;
  // 保留最近2000条
  while (log.children.length > 2000) log.removeChild(log.firstChild);
}

function escHtml(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

function updateAgentUI() {
  const list = document.getElementById('agent-list');
  list.innerHTML = '';
  Object.entries(agents).forEach(([id, online]) => {
    const div = document.createElement('div');
    div.className = 'agent-item';
    div.innerHTML = '<span class="dot ' + (online ? 'online' : '') + '"></span>' + escHtml(id);
    list.appendChild(div);
  });
}

function connect() {
  if (ws) { ws.close(); ws = null; }
  const token = document.getElementById('token-input').value.trim();
  const subscribe = document.getElementById('subscribe-input').value.trim() || '*';
  const proto = location.protocol === 'https:' ? 'wss' : 'ws';
  let url = proto + '://' + location.host + '/ws/console?subscribe=' + encodeURIComponent(subscribe);
  if (token) url += '&token=' + encodeURIComponent(token);

  appendLog('Connecting to ' + url + ' ...', 'system');
  ws = new WebSocket(url);

  ws.onopen = () => {
    document.getElementById('status').textContent = '● Connected';
    document.getElementById('status').className = 'connected';
    appendLog('Connected ✓', 'system');
    // 拉取 agent 列表
    fetchAgents();
  };

  ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data);
      if (msg.type === 'agent_status') {
        agents[msg.agentId] = msg.online;
        updateAgentUI();
        appendLog('[agent_status] ' + msg.agentId + ' → ' + (msg.online ? '🟢 online' : '🔴 offline'), 'agent_status');
      } else {
        appendLog(JSON.stringify(msg), 'log');
      }
    } catch {
      appendLog(e.data, 'log');
    }
  };

  ws.onerror = () => appendLog('WebSocket error', 'error');

  ws.onclose = () => {
    document.getElementById('status').textContent = '● Disconnected';
    document.getElementById('status').className = 'disconnected';
    appendLog('Disconnected', 'system');
  };
}

async function fetchAgents() {
  try {
    const res = await fetch('/api/agents');
    const data = await res.json();
    data.agents.forEach(a => agents[a.id] = true);
    updateAgentUI();
  } catch(e) {}
}

// 自动连接
window.onload = () => connect();
</script>
</body>
</html>`
}
