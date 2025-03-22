package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
	Meta     MetaConfig     `yaml:"meta"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port               int `yaml:"port"`
	ReadTimeoutSecs    int `yaml:"read_timeout_seconds"`
	WriteTimeoutSecs   int `yaml:"write_timeout_seconds"`
	IdleTimeoutSecs    int `yaml:"idle_timeout_seconds"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URL                     string `yaml:"url"`
	MaxConnections          int    `yaml:"max_connections"`
	MinConnections          int    `yaml:"min_connections"`
	MaxConnectionLifetimeHr int    `yaml:"max_connection_lifetime_hours"`
	MaxConnectionIdleMin    int    `yaml:"max_connection_idle_minutes"`
}

// LoggingConfig holds logging-specific configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// MetaConfig holds meta information
type MetaConfig struct {
	Version  string `yaml:"version"`
	BuildDate string `yaml:"build_date"`
	Author   string `yaml:"author"`
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(filePath string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for any missing values
	applyDefaults(&config)

	return &config, nil
}

// Load attempts to load the configuration from predefined locations
func Load() (*Config, error) {
	// Define potential config file locations
	configPaths := []string{
		"./config.yaml",
		"./config/config.yaml",
		"/etc/emqx-pg-bridge/config.yaml",
		filepath.Join(os.Getenv("HOME"), ".emqx-pg-bridge", "config.yaml"),
	}

	// Try each location
	var lastErr error
	for _, path := range configPaths {
		config, err := LoadFromFile(path)
		if err == nil {
			return config, nil
		}
		lastErr = err
	}

	// If no config file is found, create a default configuration
	config := &Config{}
	applyDefaults(config)
	
	// Log a warning that we're using default configuration
	fmt.Println("WARNING: No configuration file found. Using default configuration.")
	fmt.Println("Checked paths:", configPaths)
	
	return config, nil
}

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Server.Port)
	}

	if c.Database.URL == "" {
		return fmt.Errorf("database URL cannot be empty")
	}

	if c.Database.MaxConnections <= 0 {
		return fmt.Errorf("max connections must be positive")
	}

	return nil
}

// applyDefaults sets default values for missing configuration
func applyDefaults(config *Config) {
	// Server defaults
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeoutSecs == 0 {
		config.Server.ReadTimeoutSecs = 15  // Fixed: replaced 'a' with '15'
	}
	if config.Server.WriteTimeoutSecs == 0 {
		config.Server.WriteTimeoutSecs = 15
	}
	if config.Server.IdleTimeoutSecs == 0 {
		config.Server.IdleTimeoutSecs = 60
	}

	// Database defaults
	if config.Database.URL == "" {
		config.Database.URL = "postgres://postgres:postgres@localhost:5432/postgres"
	}
	if config.Database.MaxConnections == 0 {
		config.Database.MaxConnections = 10
	}
	if config.Database.MinConnections == 0 {
		config.Database.MinConnections = 1
	}
	if config.Database.MaxConnectionLifetimeHr == 0 {
		config.Database.MaxConnectionLifetimeHr = 1
	}
	if config.Database.MaxConnectionIdleMin == 0 {
		config.Database.MaxConnectionIdleMin = 30
	}

	// Logging defaults
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "text"
	}

	// Meta defaults
	if config.Meta.Version == "" {
		config.Meta.Version = "dev"
	}
	if config.Meta.BuildDate == "" {
		config.Meta.BuildDate = time.Now().Format("2006-01-02 15:04:05")
	}
}

// GetReadTimeout returns the read timeout as a duration
func (c *Config) GetReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeoutSecs) * time.Second
}

// GetWriteTimeout returns the write timeout as a duration
func (c *Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeoutSecs) * time.Second
}

// GetIdleTimeout returns the idle timeout as a duration
func (c *Config) GetIdleTimeout() time.Duration {
	return time.Duration(c.Server.IdleTimeoutSecs) * time.Second
}

// GetMaxConnectionLifetime returns the maximum connection lifetime as a duration
func (c *Config) GetMaxConnectionLifetime() time.Duration {
	return time.Duration(c.Database.MaxConnectionLifetimeHr) * time.Hour
}

// GetMaxConnectionIdleTime returns the maximum connection idle time as a duration
func (c *Config) GetMaxConnectionIdleTime() time.Duration {
	return time.Duration(c.Database.MaxConnectionIdleMin) * time.Minute
}