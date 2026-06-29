package main

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestMessageStorePersistsAndFiltersLogs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "messages.jsonl")
	store, err := NewMessageStore(path, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Append("host-a", []byte(`{"type":"agent_status","agentId":"host-a"}`)); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("host-a", []byte(`{"type":"log","agentId":"host-a","lines":["a"]}`)); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("host-b", []byte(`{"type":"log","agentId":"host-b","lines":["b"]}`)); err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewMessageStore(path, 10)
	if err != nil {
		t.Fatal(err)
	}
	messages := reloaded.List("host-a", 10)
	if len(messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(messages))
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(messages[0], &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope["agentId"] != "host-a" {
		t.Fatalf("unexpected envelope: %#v", envelope)
	}
}

func TestMessageStoreKeepsConfiguredLimit(t *testing.T) {
	store, err := NewMessageStore("", 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"a", "b", "c"} {
		if err := store.Append(id, []byte(`{"type":"log","lines":[]}`)); err != nil {
			t.Fatal(err)
		}
	}
	if got := len(store.List("*", 10)); got != 2 {
		t.Fatalf("got %d messages, want 2", got)
	}
}

func TestMessageStoreSupportsSeparateTrajectoryEvents(t *testing.T) {
	store, err := NewMessageStore("", 10, "trajectory")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Append("host-a", []byte(`{"type":"log"}`)); err != nil {
		t.Fatal(err)
	}
	if err := store.Append("host-a", []byte(`{"type":"trajectory","eventType":"model.completed"}`)); err != nil {
		t.Fatal(err)
	}
	if got := len(store.List("*", 10)); got != 1 {
		t.Fatalf("got %d trajectory events, want 1", got)
	}
}
