package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type trajectoryInfo struct {
	AgentID   string
	SessionID string
	Path      string
	MTime     int64
	Size      int64
}

type trajectoryEvent struct {
	Type       string                 `json:"type"`
	Timestamp  string                 `json:"ts,omitempty"`
	Provider   string                 `json:"provider,omitempty"`
	ModelID    string                 `json:"modelId,omitempty"`
	ModelAPI   string                 `json:"modelApi,omitempty"`
	TraceID    string                 `json:"traceId,omitempty"`
	RunID      string                 `json:"runId,omitempty"`
	SessionKey string                 `json:"sessionKey,omitempty"`
	Workspace  string                 `json:"workspaceDir,omitempty"`
	Seq        int64                  `json:"seq,omitempty"`
	SourceSeq  int64                  `json:"sourceSeq,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

func discoverTrajectories(root string) ([]trajectoryInfo, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var trajectories []trajectoryInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		agentID := entry.Name()
		sessionDir := filepath.Join(root, agentID, "sessions")
		files, err := os.ReadDir(sessionDir)
		if err != nil {
			continue
		}
		for _, file := range files {
			name := file.Name()
			if file.IsDir() || !strings.HasSuffix(name, ".trajectory.jsonl") {
				continue
			}
			info, err := file.Info()
			if err != nil {
				continue
			}
			trajectories = append(trajectories, trajectoryInfo{
				AgentID: agentID, SessionID: strings.TrimSuffix(name, ".trajectory.jsonl"),
				Path: filepath.Join(sessionDir, name), MTime: info.ModTime().UnixMilli(), Size: info.Size(),
			})
		}
	}
	sort.Slice(trajectories, func(i, j int) bool { return trajectories[i].MTime > trajectories[j].MTime })
	return trajectories, nil
}

func (a *agent) readNewTrajectoryLines(path string, state *fileState) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	if _, err := file.Seek(state.Offset, io.SeekStart); err != nil {
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return
	}
	state.Offset += int64(len(data))
	data = append(state.Remainder, data...)
	lastNewline := strings.LastIndexByte(string(data), '\n')
	if lastNewline < 0 {
		state.Remainder = append(state.Remainder[:0], data...)
		return
	}
	complete := data[:lastNewline]
	state.Remainder = append(state.Remainder[:0], data[lastNewline+1:]...)
	scanner := bufio.NewScanner(strings.NewReader(string(complete)))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var event trajectoryEvent
		raw := append([]byte(nil), scanner.Bytes()...)
		if json.Unmarshal(raw, &event) != nil || event.Type == "" {
			continue
		}
		payload := map[string]interface{}{
			"type": "trajectory", "protocolVersion": protocolVersion, "host": a.host,
			"openclawAgentId": state.AgentID, "sessionId": state.SessionID,
			"eventType": event.Type, "eventTimestamp": event.Timestamp,
			"provider": event.Provider, "modelId": event.ModelID, "modelApi": event.ModelAPI,
			"seq": event.Seq, "sourceSeq": event.SourceSeq, "redacted": a.cfg.TrajectoryRedact,
		}
		if a.cfg.TrajectoryRedact {
			// Redact only data (large prompt/response content); keep identifiers readable
			payload["traceId"] = event.TraceID
			payload["runId"] = event.RunID
			payload["sessionKey"] = event.SessionKey
			payload["workspaceDir"] = event.Workspace
			// Extract only usage stats from data if available
			if event.Data != nil {
				redactedData := map[string]interface{}{}
				for _, key := range []string{"usage", "promptCache", "compactionCount", "aborted", "timedOut"} {
					if v, exists := event.Data[key]; exists {
						redactedData[key] = v
					}
				}
				if len(redactedData) > 0 {
					payload["data"] = redactedData
				}
			}
		} else {
			payload["traceId"] = event.TraceID
			payload["runId"] = event.RunID
			payload["sessionKey"] = event.SessionKey
			payload["workspaceDir"] = event.Workspace
			payload["data"] = event.Data
		}
		a.send(payload)
	}
}

func hashIdentifier(value string) string {
	if value == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(sum[:8])
}
