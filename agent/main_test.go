package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverSessionsAndSummarizeAgents(t *testing.T) {
	root := t.TempDir()
	writeSession(t, root, "main", "session-a", "{\"message\":{\"role\":\"user\"}}\n")
	writeSession(t, root, "bajie", "session-b", "{\"message\":{\"role\":\"assistant\"}}\n")
	writeSession(t, root, "bajie", "ignored.trajectory", "{}\n")

	sessions, err := discoverSessions(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	agents := summarizeAgents("testdata/agents", sessions)
	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}
}

func TestReadNewLinesPreservesPartialJSONL(t *testing.T) {
	root := t.TempDir()
	path := writeSession(t, root, "main", "session-a", `{"message":{"role":"user"`)
	a := &agent{
		host:    hostInfo{ID: "host-a", Hostname: "host-a"},
		files:   make(map[string]*fileState),
		pending: make([][]byte, 0),
	}
	state := &fileState{AgentID: "main", SessionID: "session-a"}
	a.readNewLines(path, state)
	if len(a.pending) != 0 || len(state.Remainder) == 0 {
		t.Fatal("partial JSON line should be retained without sending")
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = file.WriteString("}}\n")
	_ = file.Close()
	a.readNewLines(path, state)

	if len(a.pending) != 1 {
		t.Fatalf("expected one queued log event, got %d", len(a.pending))
	}
	var message map[string]any
	if err := json.Unmarshal(a.pending[0], &message); err != nil {
		t.Fatal(err)
	}
	if message["openclawAgentId"] != "main" || message["sessionId"] != "session-a" {
		t.Fatalf("missing identity fields: %#v", message)
	}
}

func TestDisplayHostnameFallsBackToUniqueHostID(t *testing.T) {
	if got := displayHostname("localhost", "localhost-a1b2"); got != "localhost-a1b2" {
		t.Fatalf("expected host ID fallback, got %q", got)
	}
	if got := displayHostname("prod-mac", "prod-mac-a1b2"); got != "prod-mac" {
		t.Fatalf("expected hostname, got %q", got)
	}
}

func TestDiscoverTrajectories(t *testing.T) {
	root := t.TempDir()
	writeSession(t, root, "main", "session-a", "{}\n")
	writeSession(t, root, "main", "session-a.trajectory", "{}\n")

	trajectories, err := discoverTrajectories(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(trajectories) != 1 || trajectories[0].SessionID != "session-a" {
		t.Fatalf("unexpected trajectories: %#v", trajectories)
	}
}

func TestTrajectoryReportingRedactsSensitiveDataByDefault(t *testing.T) {
	root := t.TempDir()
	path := writeSession(t, root, "main", "session-a.trajectory", `{"type":"model.completed","ts":"2026-06-12T00:00:00Z","traceId":"secret-trace","runId":"secret-run","workspaceDir":"/secret/workspace","provider":"mlamp","modelId":"model-a","data":{"usage":{"input":12,"output":34,"total":46},"assistantTexts":["secret response"],"status":"ok"}}`+"\n")
	a := &agent{
		cfg:     config{TrajectoryRedact: true},
		host:    hostInfo{ID: "host-a"},
		files:   make(map[string]*fileState),
		pending: make([][]byte, 0),
	}
	a.readNewTrajectoryLines(path, &fileState{AgentID: "main", SessionID: "session-a"})
	if len(a.pending) != 1 {
		t.Fatalf("expected one trajectory event, got %d", len(a.pending))
	}
	var event map[string]interface{}
	if err := json.Unmarshal(a.pending[0], &event); err != nil {
		t.Fatal(err)
	}
	if _, ok := event["data"]; ok {
		t.Fatal("redacted trajectory event contains raw data")
	}
	if _, ok := event["workspaceDir"]; ok {
		t.Fatal("redacted trajectory event contains workspace path")
	}
	if _, ok := event["metrics"]; ok {
		t.Fatal("trajectory event contains derived metrics")
	}
}

func TestTrajectoryReportingCanIncludeRawData(t *testing.T) {
	root := t.TempDir()
	path := writeSession(t, root, "main", "session-a.trajectory", `{"type":"tool.call","workspaceDir":"/workspace","data":{"name":"exec","arguments":{"cmd":"private"}}}`+"\n")
	a := &agent{
		cfg:     config{TrajectoryRedact: false},
		host:    hostInfo{ID: "host-a"},
		files:   make(map[string]*fileState),
		pending: make([][]byte, 0),
	}
	a.readNewTrajectoryLines(path, &fileState{AgentID: "main", SessionID: "session-a"})
	var event map[string]interface{}
	if err := json.Unmarshal(a.pending[0], &event); err != nil {
		t.Fatal(err)
	}
	if event["workspaceDir"] != "/workspace" || event["data"] == nil {
		t.Fatalf("raw trajectory event missing requested data: %#v", event)
	}
	if _, ok := event["metrics"]; ok {
		t.Fatal("trajectory event contains derived metrics")
	}
}

func writeSession(t *testing.T, root, agentID, sessionID, content string) string {
	t.Helper()
	dir := filepath.Join(root, agentID, "sessions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, sessionID+".jsonl")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
