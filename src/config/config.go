package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
	"runtime"
)

type Config struct {
	Redis RedisConfig `yaml:"redis"`
	Server struct {
		Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"0.0.0.0:1321"`
	} `yaml:"server"`
}

type RedisConfig struct {
	Address string `yaml:"username" env:"REDIS_ADDRESS" env-default:"redis:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	DB int `yaml:"db" env:"REDIS_DB" env-default:"0"` // 0: use default DB
}

// Try to read variables from the config file.
// If it fails, read them from environment.
func GetConfig(filename string) (*Config, error) {
	var c Config
	path := getConfigPath(filename)
	if err := cleanenv.ReadConfig(path, &c); err != nil {
		if err := cleanenv.ReadEnv(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// Return the path on disk to the configs
func getConfigPath(configFilename string) string {
	_, currentFilename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(currentFilename), "../../configs/", configFilename)
}