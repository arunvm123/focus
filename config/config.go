package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var once sync.Once
var configuration config

type config struct {
	DomainURL                string         `yaml:"domain_url"`
	Database                 databaseConfig `yaml:"database"`
	JWTSecret                string         `yaml:"jwt_secret"`
	SendgridKey              string         `yaml:"sendgrid_key"`
	FCMServiceAccountKeyPath string         `yaml:"fcm_service_account_key_path"`
}

type databaseConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database_name"`
	Host         string `yaml:"host"`
	Port         string `json:"port"`
}

// GetConfig looks for config.yaml in the current directory and reads
// into the config struct
func GetConfig() (*config, error) {
	var err error
	once.Do(func() {
		err = cleanenv.ReadConfig("config.yaml", &configuration)
	})
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
