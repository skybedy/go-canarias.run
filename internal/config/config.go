package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                string
	DataFile            string
	EnableScraperDaemon bool
	ScraperInterval     time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:                getEnv("PORT", "8080"),
		DataFile:            getEnv("DATA_FILE", "data.json"),
		EnableScraperDaemon: getEnvBool("ENABLE_SCRAPER_DAEMON", true),
		ScraperInterval:     getEnvDurationHours("SCRAPER_INTERVAL_HOURS", 24),
	}

	if strings.TrimSpace(cfg.Port) == "" {
		return Config{}, fmt.Errorf("PORT must not be empty")
	}
	if strings.TrimSpace(cfg.DataFile) == "" {
		return Config{}, fmt.Errorf("DATA_FILE must not be empty")
	}
	if cfg.ScraperInterval <= 0 {
		return Config{}, fmt.Errorf("SCRAPER_INTERVAL_HOURS must be greater than 0")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func getEnvDurationHours(key string, fallbackHours int) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return time.Duration(fallbackHours) * time.Hour
	}
	hours, err := strconv.Atoi(raw)
	if err != nil || hours <= 0 {
		return time.Duration(fallbackHours) * time.Hour
	}
	return time.Duration(hours) * time.Hour
}
