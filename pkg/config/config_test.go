package config

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/kelseyhightower/envconfig"
)

type (
	// EnvProcessorMock is a mock that will hold the function to process the env variables
	EnvProcessorMock struct {
		err error
	}
)

// EnvProcessFuncMock is the mock function of the env variable processor
func (epm *EnvProcessorMock) EnvProcessFuncMock(prefix string, spec interface{}) error {
	return epm.err
}

// TestLoadEnvWithEnvConfigLibrary tests the load of config variables using kelseyhightower/envconfig envconfig.Process()
func TestLoadEnvWithEnvConfigLibrary(t *testing.T) {
	var tests = []struct {
		debug      bool
		slackToken string
		dbName     string
		dbUser     string
		dbPass     string
		dbHost     string
	}{
		{true, "slack-token", "db_name", "db_user", "db_pass", "db_host"},
	}

	for _, test := range tests {
		var configMock = &Config{
			Debug:      test.debug,
			SlackToken: test.slackToken,
			DB: struct {
				DatabaseName string `envconfig:"POSTGRES_DB" default:""`
				Host         string `envconfig:"POSTGRES_HOST" default:""`
				User         string `envconfig:"POSTGRES_USER" default:""`
				Password     string `envconfig:"POSTGRES_PASSWORD" default:""`
			}{
				test.dbName,
				test.dbHost,
				test.dbUser,
				test.dbPass,
			},
		}

		setEnvVariable(t, "DEBUG", strconv.FormatBool(test.debug))
		setEnvVariable(t, "SLACK_TOKEN", test.slackToken)
		setEnvVariable(t, "POSTGRES_DB", test.dbName)
		setEnvVariable(t, "POSTGRES_HOST", test.dbHost)
		setEnvVariable(t, "POSTGRES_USER", test.dbUser)
		setEnvVariable(t, "POSTGRES_PASSWORD", test.dbPass)

		c, err := LoadEnv(envconfig.Process)
		if err != nil {
			t.Fatalf("could not load config from envinronment variables: %q", err.Error())
		}

		if !reflect.DeepEqual(c, configMock) {
			t.Fatalf("could not load config from envinronment variables, got: %v, want %v", c, configMock)
		}
	}
}

// TestLoadEnvError tests if the error returned by the process function is propagated
func TestLoadEnvError(t *testing.T) {
	expectedErr := errors.New("could not load env variables")
	epm := &EnvProcessorMock{expectedErr}

	_, err := LoadEnv(epm.EnvProcessFuncMock)
	if err == nil {
		t.Fatalf("load config from envinroment should have failed")
	}
	if err != expectedErr {
		t.Fatalf("error when loading config from env variable, got: %q, want: %q", err, expectedErr)
	}
}

func setEnvVariable(t *testing.T, key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Errorf("could not set env for: key = %q, value = %q", key, value)
	}
}
