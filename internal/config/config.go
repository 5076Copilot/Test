package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddress     string
	ReadinessDelay  time.Duration
	NRFEnabled      bool
	NRFURI          string
	NodeID          string
	NotificationTTL time.Duration
	HTTPClientTO    time.Duration
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getDurationEnv(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func getBoolEnv(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

func LoadFromEnv() Config {
	return Config{
		HTTPAddress:     getenv("MB_SMF_HTTP_ADDR", ":8080"),
		ReadinessDelay:  getDurationEnv("MB_SMF_READY_DELAY", 0),
		NRFEnabled:      getBoolEnv("MB_SMF_NRF_ENABLED", false),
		NRFURI:          getenv("MB_SMF_NRF_URI", ""),
		NodeID:          getenv("MB_SMF_NODE_ID", "mbsmf-1"),
		NotificationTTL: getDurationEnv("MB_SMF_NOTIFY_TTL", 5*time.Second),
		HTTPClientTO:    getDurationEnv("MB_SMF_HTTP_CLIENT_TIMEOUT", 5*time.Second),
	}
}

