package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Config lists out configuration for all dependencies
type Config struct {
	DomainURL                string         `yaml:"domain_url"`
	Database                 DatabaseConfig `yaml:"database"`
	JWTSecret                string         `yaml:"jwt_secret"`
	SendgridKey              string         `yaml:"sendgrid_key"`
	FCMServiceAccountKeyPath string         `yaml:"fcm_service_account_key_path"`
}

// DatabaseConfig defines structure for database configuration
type DatabaseConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database_name"`
	Host         string `yaml:"host"`
	Port         string `json:"port"`
}

// GetConfig looks for config.yaml in the current directory and reads
// into the config struct
func GetConfig() (*Config, error) {
	var config Config
	err := cleanenv.ReadConfig("config.yaml", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
