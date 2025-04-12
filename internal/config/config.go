package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

// ServerConfig contains server related settings
type ServerConfig struct {
	Port              string
	RequestTimeoutSec int
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// LoggingConfig contains logging related settings
type LoggingConfig struct {
	IsDevelopment bool
	LogLevel      string
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:              "8082",
			RequestTimeoutSec: 30,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "qubool_kallyanam",
			SSLMode:  "disable",
		},
		Logging: LoggingConfig{
			IsDevelopment: true,
			LogLevel:      "info",
		},
	}
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := NewConfig()

	// Server config
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}

	if timeout := os.Getenv("SERVER_REQUEST_TIMEOUT_SEC"); timeout != "" {
		if value, err := strconv.Atoi(timeout); err == nil {
			config.Server.RequestTimeoutSec = value
		}
	}

	// Database config
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		if value, err := strconv.Atoi(port); err == nil {
			config.Database.Port = value
		}
	}

	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}

	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Logging config
	if isDev := os.Getenv("LOG_DEVELOPMENT"); isDev != "" {
		config.Logging.IsDevelopment = strings.ToLower(isDev) == "true"
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.LogLevel = logLevel
	}

	return config, nil
}
