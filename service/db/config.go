package db

import (
	"time"
)

// StdConfig standard config
func StdConfig(name string) *Config {
	return RawConfig("mysql." + name)
}

// RawConfig ...
// example: RawConfig("mysql.reader")
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	config.Name = key
	return config
}

// ClusterConfig ...
type ClusterConfig struct {
	Name string

	W *Config
	R *Config
}

// Config options
type Config struct {
	// Name for config
	Name string
	// DSN: mysql://root:secret@tcp(127.0.0.1:3307)/mysql?timeout=20s&readTimeout=20s
	DSN string `json:"dsn" toml:"dsn"`
	// Debug switch
	Debug bool `json:"debug" toml:"debug"`
	// MaxIdleConns
	MaxIdleConns int `json:"maxIdleConns" toml:"maxIdleConns"`
	// MaxOpenConns
	MaxOpenConns int `json:"maxOpenConns" toml:"maxOpenConns"`
	// ConnMaxLifetime
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" toml:"connMaxLifetime"`
	// OnDialError The error level of creating a connection. When set to panic, if the creation fails, panic immediately
	OnDialError string `json:"level" toml:"level"`

	// When recording error sql, whether to print the complete sql statement containing the parameters
	// DetailSQL bool `json:"detailSql" toml:"detailSql"`
}

// DefaultConfig return the default config
func DefaultConfig() *Config {
	return &Config{
		DSN:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Duration(300 * time.Second),
		OnDialError:     "panic",
	}
}
