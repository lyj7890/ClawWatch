package main

import (
	"encoding/json"
	"testing"
)

func TestUpdateAndInjectAgentMetadata(t *testing.T) {
	agent := &AgentConn{ID: "host-a"}
	raw := []byte(`{"type":"agent_hello","protocolVersion":2,"agentVersion":"0.2.0","host":{"hostname":"mac-a","ips":["10.0.0.2"],"os":"darwin","arch":"arm64"},"agents":[{"id":"main","sessionCount":3}]}`)

	updateAgentMetadata(raw, agent)
	if agent.Hostname != "mac-a" || agent.ProtocolVersion != 2 || len(agent.OpenClawAgents) != 1 {
		t.Fatalf("metadata not updated: %#v", agent)
	}

	enriched := injectAgentMetadata([]byte(`{"type":"log","openclawAgentId":"main","sessionId":"session-a"}`), agent)
	var message map[string]any
	if err := json.Unmarshal(enriched, &message); err != nil {
		t.Fatal(err)
	}
	if message["agentId"] != "host-a" || message["hostname"] != "mac-a" {
		t.Fatalf("metadata not injected: %#v", message)
	}
}
