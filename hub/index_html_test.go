package main

import (
	"strings"
	"testing"
)

func TestBuildIndexHTMLContainsRenderableBody(t *testing.T) {
	html := buildIndexHTML(&Config{})

	if strings.Contains(html, `<\/script>`) {
		t.Fatal("HTML contains escaped script closing tags that browsers do not recognize")
	}
	if got := strings.Count(html, "</script>"); got != 3 {
		t.Fatalf("expected 3 script closing tags, got %d", got)
	}
	if !strings.Contains(html, `<div id="app"`) {
		t.Fatal("HTML does not contain the application root")
	}
	for _, marker := range []string{"hostGroups", "selectedViewTitle", "message-list", "Search session content", "selectedSession", "sessionGroups", "incrementUnread", "unreadFor", "recentActivityWindowMs", "maxSessionsPerHost", "toggleHost", "toggleAgent", "sortedSessions", "toggleUnreadOnly", "selectedTrajectory", "/api/trajectory", "runtimeCoverage", "sessionRuntime", "summarizeRuntimeEvents", "events integrated", "Prompt context", "System prompt", "Messages submitted to model", "promptContext", "recoveredAttempts", "isFailedAssistantAttempt", "Model attempt failed", "recovered from trajectory", "raw trajectory store", "connectionStatusText", "scheduleReconnect", "Reconnecting..."} {
		if !strings.Contains(html, marker) {
			t.Fatalf("HTML does not contain redesigned UI marker %q", marker)
		}
	}
	for _, removed := range []string{"Runtime Timeline", "filteredTimeline", "viewMode='trace'", "Full event envelope"} {
		if strings.Contains(html, removed) {
			t.Fatalf("HTML still contains removed timeline marker %q", removed)
		}
	}
}
