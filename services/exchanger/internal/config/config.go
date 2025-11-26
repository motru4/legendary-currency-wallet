package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultGRPCHost = "0.0.0.0"
	defaultGRPCPort = 50051

	defaultDBHost    = "localhost"
	defaultDBPort    = 5432
	defaultDBUser    = "postgres"
	defaultDBName    = "exchanger"
	defaultDBSSLMode = "disable"

	defaultLogLevel  = "info"
	defaultLogFormat = "json"
)

type Config struct {
	GRPC GRPCConfig
	DB   DBConfig
	Log  LogConfig
}

type GRPCConfig struct {
	Host string
	Port int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type LogConfig struct {
	Level  string
	Format string
}

func New(path string) *Config {
	if path != "" {
		if err := godotenv.Load(path); err != nil {
			fmt.Fprintf(os.Stderr, "warning: cannot load env file %s: %v\n", path, err)
		}
	}

	cfg := &Config{
		GRPC: GRPCConfig{
			Host: getEnv("EXCHANGER_GRPC_HOST", defaultGRPCHost),
			Port: getEnvAsInt("EXCHANGER_GRPC_PORT", defaultGRPCPort),
		},
		DB: DBConfig{
			Host:     getEnv("EXCHANGER_DB_HOST", defaultDBHost),
			Port:     getEnvAsInt("EXCHANGER_DB_PORT", defaultDBPort),
			User:     getEnv("EXCHANGER_DB_USER", defaultDBUser),
			Password: getEnv("EXCHANGER_DB_PASSWORD", ""),
			Name:     getEnv("EXCHANGER_DB_NAME", defaultDBName),
			SSLMode:  getEnv("EXCHANGER_DB_SSLMODE", defaultDBSSLMode),
		},
		Log: LogConfig{
			Level:  getEnv("EXCHANGER_LOG_LEVEL", defaultLogLevel),
			Format: getEnv("EXCHANGER_LOG_FORMAT", defaultLogFormat),
		},
	}

	return cfg
}
