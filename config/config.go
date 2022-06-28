package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Server   Server `yaml:"server"`
	Database RDB    `yaml:"database"`
}

type Server struct {
	Port int `yaml:"port" env-default:"8080"`
}

type RDB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	TimeZone string `yaml:"time_zone"`
}

func NewConfig(configPath string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
