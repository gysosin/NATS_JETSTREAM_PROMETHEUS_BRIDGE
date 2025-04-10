package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenPort  string   `json:"listen_port"`
	NatsURL     string   `json:"nats_url"`
	Subject     string   `json:"subject"`
	AgentFilter []string `json:"agent_filter"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}
