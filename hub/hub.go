package main

import (
	"log"
	"sync"
	"time"
)

// AgentConn 已注册的 agent 连接
type AgentConn struct {
	ID       string
	Send     chan []byte
	LastSeen time.Time
}

// ConsoleConn 已订阅的 console 连接
type ConsoleConn struct {
	ID        string // 随机 id
	Subscribe string // 订阅的 agentId，"*" 表示全部
	Send      chan []byte
}

type broadcastMsg struct {
	agentID string
	data    []byte
}

// Hub 核心结构
type Hub struct {
	mu       sync.RWMutex
	agents   map[string]*AgentConn
	consoles map[string]*ConsoleConn

	registerAgent   chan *AgentConn
	unregisterAgent chan *AgentConn

	registerConsole   chan *ConsoleConn
	unregisterConsole chan *ConsoleConn

	broadcast chan *broadcastMsg
}

func NewHub() *Hub {
	return &Hub{
		agents:            make(map[string]*AgentConn),
		consoles:          make(map[string]*ConsoleConn),
		registerAgent:     make(chan *AgentConn, 16),
		unregisterAgent:   make(chan *AgentConn, 16),
		registerConsole:   make(chan *ConsoleConn, 16),
		unregisterConsole: make(chan *ConsoleConn, 16),
		broadcast:         make(chan *broadcastMsg, 256),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case agent := <-h.registerAgent:
			h.mu.Lock()
			h.agents[agent.ID] = agent
			count := len(h.agents)
			h.mu.Unlock()
			log.Printf("[hub] agent registered: %s (total: %d)", agent.ID, count)
			h.notifyConsolesStatus(agent.ID, true)

		case agent := <-h.unregisterAgent:
			h.mu.Lock()
			if _, ok := h.agents[agent.ID]; ok {
				delete(h.agents, agent.ID)
				close(agent.Send)
			}
			count := len(h.agents)
			h.mu.Unlock()
			log.Printf("[hub] agent disconnected: %s (total: %d)", agent.ID, count)
			h.notifyConsolesStatus(agent.ID, false)

		case console := <-h.registerConsole:
			h.mu.Lock()
			h.consoles[console.ID] = console
			h.mu.Unlock()
			log.Printf("[hub] console connected: %s subscribe=%s", console.ID, console.Subscribe)
			go h.sendCurrentAgentsStatus(console)

		case console := <-h.unregisterConsole:
			h.mu.Lock()
			if _, ok := h.consoles[console.ID]; ok {
				delete(h.consoles, console.ID)
				close(console.Send)
			}
			h.mu.Unlock()
			log.Printf("[hub] console disconnected: %s", console.ID)

		case msg := <-h.broadcast:
			h.fanout(msg)

		case <-ticker.C:
			h.logStats()
		}
	}
}

func (h *Hub) fanout(msg *broadcastMsg) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.consoles {
		if c.Subscribe == "*" || c.Subscribe == msg.agentID {
			select {
			case c.Send <- msg.data:
			default:
				log.Printf("[hub] console %s send buffer full, dropping", c.ID)
			}
		}
	}
}

func (h *Hub) notifyConsolesStatus(agentID string, online bool) {
	data := buildAgentStatusJSON(agentID, online)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.consoles {
		if c.Subscribe == "*" || c.Subscribe == agentID {
			select {
			case c.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) sendCurrentAgentsStatus(c *ConsoleConn) {
	h.mu.RLock()
	ids := make([]string, 0, len(h.agents))
	for id := range h.agents {
		ids = append(ids, id)
	}
	h.mu.RUnlock()

	for _, id := range ids {
		if c.Subscribe == "*" || c.Subscribe == id {
			data := buildAgentStatusJSON(id, true)
			select {
			case c.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) logStats() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	log.Printf("[hub] stats: agents=%d consoles=%d", len(h.agents), len(h.consoles))
}
