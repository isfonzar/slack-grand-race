package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Specification represents structured configuration variables
type Specification struct {
	Debug      bool   `envconfig:"DEBUG" default:"false"`
	SlackToken string `envconfig:"SLACK_TOKEN" default:""`
	DB         struct {
		DatabaseName string `envconfig:"POSTGRES_DB" default:""`
		Host         string `envconfig:"POSTGRES_HOST" default:""`
		User         string `envconfig:"POSTGRES_USER" default:""`
		Password     string `envconfig:"POSTGRES_PASSWORD" default:""`
	}
}

// LoadEnv loads config variables into Specification
func LoadEnv() (*Specification, error) {
	var conf Specification
	err := envconfig.Process("", &conf)
	if err != nil {
		return nil, errors.Wrap(err, "error loading configs from env vars")
	}

	return &conf, nil
}
