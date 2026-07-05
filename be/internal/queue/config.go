package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Connection ConnectionConfig          `json:"connection"`
	Streams    map[string]StreamOptions  `json:"streams"`
}

type ConnectionConfig struct {
	Enabled       bool   `json:"enabled"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	MaxReconnects int    `json:"max_reconnects"`
	ReconnectWait string `json:"reconnect_wait"`
}

type StreamOptions struct {
	Retention       string                      `json:"retention"`
	Storage         string                      `json:"storage"`
	MaxAge          string                      `json:"max_age"`
	MaxMsgs         int64                       `json:"max_msgs"`
	MaxBytes        int64                       `json:"max_bytes"`
	Discard         string                      `json:"discard"`
	DuplicateWindow string                      `json:"duplicate_window"`
	Consumers       map[string]ConsumerOptions  `json:"consumers"`
}

type ConsumerOptions struct {
	AckPolicy     string `json:"ack_policy"`
	AckWait       string `json:"ack_wait"`
	MaxDeliver    int    `json:"max_deliver"`
	MaxAckPending int    `json:"max_ack_pending"`
	DeliverPolicy string `json:"deliver_policy"`
	ReplayPolicy  string `json:"replay_policy"`
}

type StreamDefinition struct {
	Name     string
	Subjects []string
}

type ConsumerDefinition struct {
	Name          string
	FilterSubject string
	Handler       string
}

func LoadConfig() Config {
	path := os.Getenv("QUEUE_CONFIG_FILE")
	if path == "" {
		path = "internal/queue/nats.json"
	}

	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		applyEnvOverrides(&cfg)
		return cfg
	}

	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		fmt.Fprintf(os.Stderr, "queue: failed to parse %s: %v\n", path, err)
		applyEnvOverrides(&cfg)
		return cfg
	}

	mergeConfig(&cfg, fileCfg)
	applyEnvOverrides(&cfg)
	return cfg
}

func defaultConfig() Config {
	return Config{
		Connection: ConnectionConfig{
			Enabled:       false,
			URL:           "nats://localhost:4222",
			Name:          "system-be",
			MaxReconnects: -1,
			ReconnectWait: "2s",
		},
		Streams: map[string]StreamOptions{
			StreamKeySearch: {
				Retention:       "limits",
				Storage:         "file",
				MaxAge:          "168h",
				MaxMsgs:         100000,
				MaxBytes:        1073741824,
				Discard:         "old",
				DuplicateWindow: "2m",
				Consumers: map[string]ConsumerOptions{
					ConsumerKeySearchOutbox: {
						AckPolicy:     "explicit",
						AckWait:       "30s",
						MaxDeliver:    5,
						MaxAckPending: 64,
						DeliverPolicy: "all",
						ReplayPolicy:  "instant",
					},
				},
			},
		},
	}
}

func mergeConfig(base *Config, fileCfg Config) {
	if fileCfg.Connection.URL != "" || fileCfg.Connection.Enabled || fileCfg.Connection.Name != "" {
		base.Connection = fileCfg.Connection
	}
	if fileCfg.Connection.MaxReconnects != 0 {
		base.Connection.MaxReconnects = fileCfg.Connection.MaxReconnects
	}
	if fileCfg.Connection.ReconnectWait != "" {
		base.Connection.ReconnectWait = fileCfg.Connection.ReconnectWait
	}
	if len(fileCfg.Streams) > 0 {
		base.Streams = fileCfg.Streams
	}
}

func applyEnvOverrides(cfg *Config) {
	if raw := strings.TrimSpace(os.Getenv("NATS_ENABLED")); raw != "" {
		cfg.Connection.Enabled = parseBool(raw)
	}
	if url := strings.TrimSpace(os.Getenv("NATS_URL")); url != "" {
		cfg.Connection.URL = url
	}
}

func resolveStream(key string) (StreamDefinition, error) {
	switch key {
	case StreamKeySearch:
		return StreamDefinition{
			Name:     StreamSearch,
			Subjects: []string{SubjectSearchWildcard},
		}, nil
	default:
		return StreamDefinition{}, fmt.Errorf("queue: unknown stream key %q", key)
	}
}

func resolveConsumer(streamKey, consumerKey string) (ConsumerDefinition, error) {
	switch streamKey {
	case StreamKeySearch:
		switch consumerKey {
		case ConsumerKeySearchOutbox:
			return ConsumerDefinition{
				Name:          ConsumerSearchOutbox,
				FilterSubject: SubjectSearchOutbox,
				Handler:       HandlerSearchOutbox,
			}, nil
		}
	}
	return ConsumerDefinition{}, fmt.Errorf("queue: unknown consumer key %q for stream %q", consumerKey, streamKey)
}

func parseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
