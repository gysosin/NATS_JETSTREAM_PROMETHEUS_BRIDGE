package server

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type MetricsPayload struct {
	SystemName string `json:"system_name"`
	Metrics    string `json:"metrics"`
}

// CacheEntry holds metrics + last update timestamp
type CacheEntry struct {
	Metrics     string
	LastUpdated int64 // Unix timestamp
}

// Collector manages metric cache per agent
type Collector struct {
	Cache     map[string]CacheEntry
	Lock      sync.RWMutex
	TTL       int64         // how long (seconds) before metrics expire
	CleanupHz time.Duration // how often to clean expired entries
}

// Start sets up the NATS subscription and cleanup worker
func (c *Collector) Start(natsURL, subject string) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return err
	}
	js, _ := nc.JetStream()

	_, err = js.Subscribe(subject, func(msg *nats.Msg) {
		var payload MetricsPayload
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Println("Bad payload:", err)
			return
		}

		// Inject system_name label into each metric line
		normalized := injectSystemNameLabel(payload.Metrics, payload.SystemName)

		c.Lock.Lock()
		c.Cache[payload.SystemName] = CacheEntry{
			Metrics:     normalized,
			LastUpdated: time.Now().Unix(),
		}
		c.Lock.Unlock()
	})
	if err != nil {
		return err
	}

	// Start background cache cleaner
	go c.cleanupExpired()
	return nil
}

// cleanupExpired removes cache entries not updated in TTL seconds
func (c *Collector) cleanupExpired() {
	ticker := time.NewTicker(c.CleanupHz)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()

		c.Lock.Lock()
		for key, entry := range c.Cache {
			if now-entry.LastUpdated > c.TTL {
				delete(c.Cache, key)
				log.Printf("Evicted stale metrics from: %s", key)
			}
		}
		c.Lock.Unlock()
	}
}

// injectSystemNameLabel rewrites metrics to include system_name label
func injectSystemNameLabel(metrics, systemName string) string {
	lines := strings.Split(metrics, "\n")
	var output []string

	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			output = append(output, line)
			continue
		}

		if strings.Contains(line, "{") {
			line = strings.Replace(line, "{", "{system_name=\""+systemName+"\",", 1)
		} else {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				line = parts[0] + "{system_name=\"" + systemName + "\"} " + strings.Join(parts[1:], " ")
			}
		}
		output = append(output, line)
	}
	return strings.Join(output, "\n")
}
