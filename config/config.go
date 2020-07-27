package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

var configuration config

type config struct {
	Port                     string         `yaml:"port" env:"PORT"`
	DomainURL                string         `yaml:"domain_url" env:"DOMAIN_URL"`
	Database                 databaseConfig `yaml:"database" env:"DATABASE"`
	JWTSecret                string         `yaml:"jwt_secret" env:"JWT_SECRET"`
	SendgridKey              string         `yaml:"sendgrid_key" env:"SENDGRID_KEY"`
	EmailTemplate            emailTemplate  `yaml:"email_template" env:"EMAIL_TEMPLATE"`
	FCMServiceAccountKeyPath string         `yaml:"fcm_service_account_key_path" env:"FCM_SERVICE_ACCOUNT_KEY_PATH"`
	AdminIDs                 adminIDs       `yaml:"admin_ids" env:"ADMIN_IDS"`
}

type databaseConfig struct {
	User         string `yaml:"user" env:"DB_USER"`
	Password     string `yaml:"password" env:"DB_PASSWORD"`
	DatabaseName string `yaml:"database_name" env:"DB_NAME"`
	Host         string `yaml:"host" env:"DB_HOST"`
	Port         string `json:"port" env:"DB_PORT"`
}

// Specifies ID for each template
type emailTemplate struct {
	EmailValidation    string `yaml:"email_validation" env:"EMAIL_VALIDATION"`
	ForgotPassword     string `yaml:"forgot_password" env:"FORGOT_PASSWORD"`
	OrganisationInvite string `yaml:"organisation_invite" env:"ORGANISATION_INVITE"`
}

type adminIDs []int

func Initialise(filepath string, env bool) (*config, error) {
	var err error

	if env {
		err = cleanenv.ReadEnv(&configuration)
	} else {
		err = cleanenv.ReadConfig(filepath, &configuration)
	}

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

func (ids *adminIDs) SetValue(s string) error {
	if s == "" {
		return fmt.Errorf("field value can't be empty")
	}

	stringIDs := strings.Split(s, " ")

	for i := 0; i < len(stringIDs); i++ {
		v, err := strconv.Atoi(stringIDs[i])
		if err != nil {
			return fmt.Errorf("Provide valid values")
		}

		*ids = append(*ids, v)
	}
	return nil
}
