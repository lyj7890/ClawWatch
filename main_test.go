package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestTrajectoryHistoryRemovesLegacyMetrics(t *testing.T) {
	store, err := NewMessageStore("", 10, "trajectory")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Append("host-a", []byte(`{"type":"trajectory","eventType":"tool.call","data":{"name":"exec"},"metrics":{"toolCount":1}}`)); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest("GET", "/api/trajectory?agentId=host-a", nil)
	response := httptest.NewRecorder()
	tokens, err := NewTokenStore("")
	if err != nil {
		t.Fatal(err)
	}
	hub := NewHub(nil, store, tokens)
	cfg := &Config{}
	handleTrajectoryHistory(hub, cfg, store)(response, request)

	var body struct {
		Messages []map[string]interface{} `json:"messages"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Messages) != 1 {
		t.Fatalf("got %d trajectory events, want 1", len(body.Messages))
	}
	if _, ok := body.Messages[0]["metrics"]; ok {
		t.Fatal("trajectory API returned legacy metrics")
	}
	if body.Messages[0]["data"] == nil {
		t.Fatal("trajectory API removed original event data")
	}
}
