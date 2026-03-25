package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	App    AppConfig
	Server ServerConfig
}

// AppConfig contains application-level settings.
type AppConfig struct {
	Name    string
	Env     string
	Version string
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// Addr returns the host:port address string.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT value: %w", err)
	}

	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "myapp"),
			Env:     getEnv("APP_ENV", "development"),
			Version: getEnv("APP_VERSION", "0.1.0"),
		},
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            port,
			ReadTimeout:     getDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}
