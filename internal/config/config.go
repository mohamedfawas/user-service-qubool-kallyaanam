package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
	JWT      JWTConfig
}

// ServerConfig contains server related settings
type ServerConfig struct {
	Port              string
	RequestTimeoutSec int
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	DBName        string
	SSLMode       string
	RunMigrations bool
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

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

func validateConfig(config *Config) error {
	// Validate JWT configuration
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// Validate database configuration
	if config.Database.User == "" || config.Database.Password == "" {
		return fmt.Errorf("database credentials (DB_USER, DB_PASSWORD) are required")
	}

	return nil
}

// LoadConfig loads configuration using Viper
func LoadConfig() (*Config, error) {
	// Initialize viper
	v := viper.New()

	// Set up viper to read environment variables
	v.AutomaticEnv()

	// Set default values
	setDefaults(v)

	// Create and return the config
	config := &Config{
		Server: ServerConfig{
			Port:              v.GetString("SERVER_PORT"),
			RequestTimeoutSec: v.GetInt("SERVER_REQUEST_TIMEOUT_SEC"),
		},
		Database: DatabaseConfig{
			Host:          v.GetString("DB_HOST"),
			Port:          v.GetInt("DB_PORT"),
			User:          v.GetString("DB_USER"),
			Password:      v.GetString("DB_PASSWORD"),
			DBName:        v.GetString("DB_NAME"),
			SSLMode:       v.GetString("DB_SSL_MODE"),
			RunMigrations: v.GetBool("DB_RUN_MIGRATIONS"),
		},
		Logging: LoggingConfig{
			IsDevelopment: v.GetBool("LOG_DEVELOPMENT"),
			LogLevel:      v.GetString("LOG_LEVEL"),
		},
		JWT: JWTConfig{
			Secret:        v.GetString("JWT_SECRET"),
			TokenExpiry:   v.GetDuration("JWT_TOKEN_EXPIRY"),
			RefreshExpiry: v.GetDuration("JWT_REFRESH_EXPIRY"),
			Issuer:        v.GetString("JWT_ISSUER"),
		},
	}

	// Add this before returning:
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// setDefaults configures all the default values for viper
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("SERVER_PORT", "8082")
	v.SetDefault("SERVER_REQUEST_TIMEOUT_SEC", 30)

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "user_db")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_RUN_MIGRATIONS", true)

	// Logging defaults
	v.SetDefault("LOG_DEVELOPMENT", true)
	v.SetDefault("LOG_LEVEL", "info")

	// JWT defaults - using the same configuration across services
	v.SetDefault("JWT_TOKEN_EXPIRY", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRY", "24h")
	v.SetDefault("JWT_ISSUER", "qubool-kallyaanam-api")
}

// NewConfig creates a new configuration with default values - kept for backward compatibility
func NewConfig() *Config {
	config, _ := LoadConfig()
	return config
}
