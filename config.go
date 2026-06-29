package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	AgentToken   string
	ConsoleToken string
	StoragePath    string
	HistoryLimit   int
	MetricsPath    string
	MetricsLimit   int
	TokenStorePath string
}

func loadConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "port", getEnv("HUB_PORT", "4848"), "Hub listen port")
	flag.StringVar(&cfg.AgentToken, "agent-token", getEnv("HUB_AGENT_TOKEN", ""), "Token required for agent registration")
	flag.StringVar(&cfg.ConsoleToken, "console-token", getEnv("HUB_CONSOLE_TOKEN", ""), "Token required for console access")
	flag.StringVar(&cfg.StoragePath, "storage-path", getEnv("HUB_STORAGE_PATH", "data/messages.jsonl"), "Path for persisted log messages")
	flag.IntVar(&cfg.HistoryLimit, "history-limit", getEnvInt("HUB_HISTORY_LIMIT", 2000), "Maximum persisted log envelopes")
	flag.StringVar(&cfg.MetricsPath, "metrics-path", getEnv("HUB_METRICS_PATH", "data/metrics.jsonl"), "Path for persisted trajectory metrics")
	flag.IntVar(&cfg.MetricsLimit, "metrics-limit", getEnvInt("HUB_METRICS_LIMIT", 10000), "Maximum persisted trajectory metric events")
	flag.StringVar(&cfg.TokenStorePath, "token-store-path", getEnv("HUB_TOKEN_STORE_PATH", "data/tokens.json"), "Path for persisted per-agent console tokens")
	flag.Parse()

	return cfg
}

func getEnvInt(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
