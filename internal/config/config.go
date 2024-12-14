package config

import (
	"fmt"
	"time"

	"github.com/oleksiip-aiola/erdtree/internal/db"
	"github.com/spf13/viper"
)

type MasterServerConfig struct {
	Server         ServerConfig
	Database       *db.Config
	WAL            WALConfig
	SlaveAddresses []string
	SyncInterval   time.Duration
	BatchSize      int
}
type SlaveServerConfig struct {
	Server   ServerConfig
	Database *db.Config
	WAL      WALConfig
}
type Config struct {
	Master MasterServerConfig
	Slave  SlaveServerConfig
	Env    string
}

type ServerConfig struct {
	Port int
	Host string
}

type WALConfig struct {
	Directory    string
	SyncInterval time.Duration
}

type SlaveConfig struct {
	MasterAddress string
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found. Using default values and environment variables.")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 50051)
	viper.SetDefault("server.host", "0.0.0.0")

	// Database defaults
	viper.SetDefault("database.maxsize", 1000000)
	viper.SetDefault("database.gcinterval", "1m")

	// WAL defaults
	viper.SetDefault("wal.directory", "/var/lib/kvstore/wal")
	viper.SetDefault("wal.syncinterval", "100ms")

	// Master defaults
	viper.SetDefault("master.syncinterval", "10s")
	viper.SetDefault("master.batchsize", 100)

	// Slave defaults
	// No default for master address, as it's required for slave mode
}

func (c *Config) IsMaster() bool {
	return len(c.Master.SlaveAddresses) > 0
}

func ValidateConfig(config *Config) error {
	if config.Master.Server.Port <= 0 || config.Master.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Master.Server.Port)
	}

	if config.Master.Database.MaxSize <= 0 {
		return fmt.Errorf("invalid database max size: %d", config.Master.Database.MaxSize)
	}

	if config.Master.Database.GCInterval <= 0 {
		return fmt.Errorf("invalid database GC interval: %v", config.Master.Database.GCInterval)
	}

	if config.Master.WAL.SyncInterval <= 0 {
		return fmt.Errorf("invalid WAL sync interval: %v", config.Master.WAL.SyncInterval)
	}

	if config.IsMaster() {
		if config.Master.SyncInterval <= 0 {
			return fmt.Errorf("invalid master sync interval: %v", config.Master.SyncInterval)
		}
		if config.Master.BatchSize <= 0 {
			return fmt.Errorf("invalid master batch size: %d", config.Master.BatchSize)
		}
	}
	// else {
	// 	if config.Master.Slave.MasterAddress == "" {
	// 		return fmt.Errorf("master address is required for slave mode")
	// 	}
	// }

	return nil
}
