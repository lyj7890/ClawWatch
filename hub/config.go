package main

import (
	"flag"
	"os"
)

type Config struct {
	Port        string
	AgentToken  string
	ConsoleToken string
}

func loadConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Port, "port", getEnv("HUB_PORT", "4848"), "Hub listen port")
	flag.StringVar(&cfg.AgentToken, "agent-token", getEnv("HUB_AGENT_TOKEN", ""), "Token required for agent registration")
	flag.StringVar(&cfg.ConsoleToken, "console-token", getEnv("HUB_CONSOLE_TOKEN", ""), "Token required for console access")
	flag.Parse()

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
