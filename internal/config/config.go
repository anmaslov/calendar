package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Exchange ExchangeConfig `yaml:"exchange"`
	Sync     SyncConfig     `yaml:"sync"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	SSLMode      string `yaml:"ssl_mode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// ExchangeConfig holds Exchange server configuration.
type ExchangeConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
}

// SyncConfig holds synchronization configuration.
type SyncConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
	// SyncDays defines how many days ahead to sync events
	SyncDays int `yaml:"sync_days"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DSN returns the database connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Load loads configuration from YAML file.
func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	// Set defaults
	cfg.setDefaults()

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables for sensitive data
	cfg.overrideFromEnv()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) setDefaults() {
	c.Server.Port = 8080
	c.Server.ReadTimeout = 15 * time.Second
	c.Server.WriteTimeout = 15 * time.Second
	c.Server.IdleTimeout = 60 * time.Second
	c.Server.ShutdownTimeout = 30 * time.Second

	c.Database.Host = "localhost"
	c.Database.Port = 5432
	c.Database.User = "calendar"
	c.Database.Name = "calendar"
	c.Database.SSLMode = "disable"
	c.Database.MaxOpenConns = 25
	c.Database.MaxIdleConns = 5

	c.Logging.Level = "info"
	c.Logging.Format = "json"

	c.Sync.Enabled = false
	c.Sync.Interval = 5 * time.Minute
	c.Sync.SyncDays = 30
}

// overrideFromEnv allows overriding sensitive values from environment variables.
func (c *Config) overrideFromEnv() {
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.Database.Password = v
	}
	if v := os.Getenv("EXCHANGE_PASSWORD"); v != "" {
		c.Exchange.Password = v
	}
}

func (c *Config) validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("database password is required (set in config or DB_PASSWORD env)")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	return nil
}
