package mq

import (
	"fmt"
	"strings"
	"time"

	"github.com/stellhub/stellar/config"
)

func durationFromConfig(value string, name string, fallback time.Duration) (time.Duration, error) {
	if strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("stellar: invalid mq %s %q: %w", name, value, err)
	}
	return duration, nil
}

func cloneProperties(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func mergeProperties(values ...map[string]string) map[string]string {
	merged := map[string]string{}
	for _, value := range values {
		for key, item := range value {
			if strings.TrimSpace(key) == "" {
				continue
			}
			merged[key] = item
		}
	}
	if len(merged) == 0 {
		return nil
	}
	return merged
}

func brokersFromConfig(cfg *config.MQConfig) []string {
	if cfg == nil || len(cfg.Brokers) == 0 {
		return []string{"localhost:9092"}
	}
	brokers := make([]string, 0, len(cfg.Brokers))
	for _, broker := range cfg.Brokers {
		if strings.TrimSpace(broker) == "" {
			continue
		}
		brokers = append(brokers, strings.TrimSpace(broker))
	}
	if len(brokers) == 0 {
		return []string{"localhost:9092"}
	}
	return brokers
}
