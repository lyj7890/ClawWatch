package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// TokenStore 维护 agentId -> consoleToken 的映射，并持久化到 JSON 文件。
// 一个 token 对应一个 agentId（一机一 token）。
type TokenStore struct {
	mu     sync.RWMutex
	path   string
	tokens map[string]string // agentId -> token
}

// NewTokenStore 加载或初始化 token 存储。path 为空表示仅内存模式。
func NewTokenStore(path string) (*TokenStore, error) {
	s := &TokenStore{
		path:   path,
		tokens: make(map[string]string),
	}
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

func (s *TokenStore) load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	m := make(map[string]string)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	s.tokens = m
	return nil
}

// saveLocked 在持有写锁时写入磁盘（原子替换）。
func (s *TokenStore) saveLocked() error {
	if s.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(s.tokens, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// GetOrCreate 返回 agentId 对应的 token，不存在则生成一个并持久化。
func (s *TokenStore) GetOrCreate(agentID string) (string, error) {
	s.mu.RLock()
	if tok, ok := s.tokens[agentID]; ok {
		s.mu.RUnlock()
		return tok, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	// double-check
	if tok, ok := s.tokens[agentID]; ok {
		return tok, nil
	}
	tok, err := generateToken()
	if err != nil {
		return "", err
	}
	s.tokens[agentID] = tok
	if err := s.saveLocked(); err != nil {
		// 回滚内存以避免内存与磁盘不一致
		delete(s.tokens, agentID)
		return "", err
	}
	return tok, nil
}

// Get 返回 agentId 对应的 token（若存在）。
func (s *TokenStore) Get(agentID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tok, ok := s.tokens[agentID]
	return tok, ok
}

// AgentIDForToken 反查 token 绑定的 agentId。
func (s *TokenStore) AgentIDForToken(token string) (string, bool) {
	if token == "" {
		return "", false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for agentID, tok := range s.tokens {
		if tok == token {
			return agentID, true
		}
	}
	return "", false
}

// generateToken 生成一个 32 字符（16 字节）随机十六进制 token。
func generateToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
