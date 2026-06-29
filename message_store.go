package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type storedMessage struct {
	AgentID string          `json:"agentId"`
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
}

// MessageStore keeps a bounded in-memory history backed by an append-only file.
type MessageStore struct {
	mu       sync.RWMutex
	path     string
	maxItems int
	items    []storedMessage
	appends  int
	types    map[string]bool
}

func NewMessageStore(path string, maxItems int, messageTypes ...string) (*MessageStore, error) {
	if maxItems <= 0 {
		maxItems = 2000
	}
	types := make(map[string]bool)
	for _, messageType := range messageTypes {
		types[messageType] = true
	}
	if len(types) == 0 {
		types["log"] = true
	}
	s := &MessageStore{path: path, maxItems: maxItems, types: types}
	if path == "" {
		return s, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *MessageStore) load() error {
	f, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		var item storedMessage
		if json.Unmarshal(scanner.Bytes(), &item) == nil && s.types[item.Type] && len(item.Data) > 0 {
			s.items = append(s.items, item)
			if len(s.items) > s.maxItems {
				s.items = s.items[len(s.items)-s.maxItems:]
			}
		}
	}
	return scanner.Err()
}

func (s *MessageStore) Append(agentID string, data []byte) error {
	var envelope struct {
		Type string `json:"type"`
	}
	if json.Unmarshal(data, &envelope) != nil || !s.types[envelope.Type] {
		return nil
	}

	item := storedMessage{
		AgentID: agentID,
		Type:    envelope.Type,
		Data:    append(json.RawMessage(nil), data...),
	}
	line, err := json.Marshal(item)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
	if len(s.items) > s.maxItems {
		s.items = s.items[len(s.items)-s.maxItems:]
	}
	if s.path == "" {
		return nil
	}

	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(append(line, '\n'))
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}

	s.appends++
	if s.appends >= s.maxItems {
		s.appends = 0
		return s.compactLocked()
	}
	return nil
}

func (s *MessageStore) List(agentID string, limit int) []json.RawMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > s.maxItems {
		limit = s.maxItems
	}
	result := make([]json.RawMessage, 0, limit)
	for i := len(s.items) - 1; i >= 0 && len(result) < limit; i-- {
		item := s.items[i]
		if agentID == "*" || agentID == "" || item.AgentID == agentID {
			result = append(result, append(json.RawMessage(nil), item.Data...))
		}
	}
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func (s *MessageStore) compactLocked() error {
	tmp := s.path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	for _, item := range s.items {
		if err := enc.Encode(item); err != nil {
			f.Close()
			return err
		}
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
