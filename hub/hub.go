package main

import (
	"log"
	"sync"
	"time"
)

// SessionInfo 表示一个 OpenClaw session
type SessionInfo struct {
	AgentID   string `json:"agentId"`
	SessionID string `json:"sessionId"`
	MTime     int64  `json:"mtime"`
	Size      int64  `json:"size"`
}

// AgentConn 已注册的 agent 连接
type AgentConn struct {
	ID              string
	Hostname        string
	HostIPs         []string
	OS              string
	Arch            string
	AgentVersion    string
	ProtocolVersion int
	OpenClawAgents  []OpenClawAgentInfo
	Sessions        []SessionInfo
	Send            chan []byte
	LastSeen        time.Time
}

type OpenClawAgentInfo struct {
	ID           string `json:"id"`
	SessionCount int    `json:"sessionCount"`
	LastActivity int64  `json:"lastActivity,omitempty"`
}

// ConsoleConn 已订阅的 console 连接
type ConsoleConn struct {
	ID        string // 随机 id
	Subscribe string // 订阅的 agentId，"*" 表示全部
	Send      chan []byte

	// AllowedAgentID 是该 console 凭 token 解析出的唯一可见 agentId。
	// 为空且 Admin=true 时表示可见全部。
	AllowedAgentID string
	// Admin 表示该 console 使用了全局 ConsoleToken（或未配置 token 的开放模式），可见全部 agent。
	Admin bool
}

// canSee 判断该 console 是否有权查看指定 agentId 的数据。
func (c *ConsoleConn) canSee(agentID string) bool {
	if c.Admin {
		return true
	}
	return c.AllowedAgentID != "" && c.AllowedAgentID == agentID
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
	store     *MessageStore
	metrics   *MessageStore
	tokens    *TokenStore

	// pendingRequests 追踪等待响应的请求
	pendingMu       sync.Mutex
	pendingRequests map[string]chan []byte
}

func NewHub(store, metrics *MessageStore, tokens *TokenStore) *Hub {
	return &Hub{
		agents:            make(map[string]*AgentConn),
		consoles:          make(map[string]*ConsoleConn),
		registerAgent:     make(chan *AgentConn, 16),
		unregisterAgent:   make(chan *AgentConn, 16),
		registerConsole:   make(chan *ConsoleConn, 16),
		unregisterConsole: make(chan *ConsoleConn, 16),
		broadcast:         make(chan *broadcastMsg, 256),
		store:             store,
		metrics:           metrics,
		tokens:            tokens,
		pendingRequests:   make(map[string]chan []byte),
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
			if current, ok := h.agents[agent.ID]; ok && current == agent {
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
			if h.store != nil {
				if err := h.store.Append(msg.agentID, msg.data); err != nil {
					log.Printf("[hub] persist message: %v", err)
				}
			}
			if h.metrics != nil {
				if err := h.metrics.Append(msg.agentID, msg.data); err != nil {
					log.Printf("[hub] persist metrics: %v", err)
				}
			}
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
		if !c.canSee(msg.agentID) {
			continue
		}
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
	h.mu.RLock()
	defer h.mu.RUnlock()
	agent := h.agents[agentID]
	if agent == nil {
		agent = &AgentConn{ID: agentID}
	}
	data := buildAgentStatusJSON(agent, online)
	for _, c := range h.consoles {
		if !c.canSee(agentID) {
			continue
		}
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
	agents := make([]*AgentConn, 0, len(h.agents))
	for _, agent := range h.agents {
		agents = append(agents, agent)
	}
	h.mu.RUnlock()

	for _, agent := range agents {
		if !c.canSee(agent.ID) {
			continue
		}
		if c.Subscribe == "*" || c.Subscribe == agent.ID {
			data := buildAgentStatusJSON(agent, true)
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
