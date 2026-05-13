package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// buildAgentStatusJSON 构造 agent 状态通知的 JSON
func buildAgentStatusJSON(agentID string, online bool) []byte {
	msg := map[string]interface{}{
		"type":      "agent_status",
		"agentId":   agentID,
		"online":    online,
		"timestamp": time.Now().UnixMilli(),
	}
	data, _ := json.Marshal(msg)
	return data
}

// generateID 生成简单唯一 ID
func generateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
