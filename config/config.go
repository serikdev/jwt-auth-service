package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerCfg
	Database DBCfg
	Logger   LoggerCfg
}

type ServerCfg struct {
	Port      string
	Host      string
	ReadTime  time.Duration
	WriteTime time.Duration
}

type DBCfg struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLmode  string
}

type LoggerCfg struct {
	Level string
}

func LoadCfg() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed load env file: %w", err)
	}
	readTimeout, err := time.ParseDuration(getEnv("SERVER_READ_TIME", "60s"))
	if err != nil {
		return nil, fmt.Errorf("parsing read timeout: %w", err)
	}
	writeTimeout, err := time.ParseDuration(getEnv("SERVER_WRITE_TIME", "60s"))
	if err != nil {
		return nil, fmt.Errorf("parsing write timeout: %w", err)
	}

	return &Config{
		Server: ServerCfg{
			Port:      getEnv("SERVER_PORT", "8080"),
			Host:      getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTime:  readTimeout,
			WriteTime: writeTimeout,
		},
		Database: DBCfg{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "jwt-service"),
			SSLmode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Logger: LoggerCfg{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
