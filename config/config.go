package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

var configuration config

type config struct {
	Port                     string         `yaml:"port"`
	DomainURL                string         `yaml:"domain_url"`
	Database                 databaseConfig `yaml:"database"`
	JWTSecret                string         `yaml:"jwt_secret"`
	SendgridKey              string         `yaml:"sendgrid_key"`
	FCMServiceAccountKeyPath string         `yaml:"fcm_service_account_key_path"`
	AdminIDs                 []int          `yaml:"admin_ids"`
}

type databaseConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database_name"`
	Host         string `yaml:"host"`
	Port         string `json:"port"`
}

func Initialise(filepath string) (*config, error) {
	err := cleanenv.ReadConfig(filepath, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

// GetConfig looks for config.yaml in the current directory and reads
// into the config struct
func GetConfig() (*config, error) {
	return &configuration, nil
}
