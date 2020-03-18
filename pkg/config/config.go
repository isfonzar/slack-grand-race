package config

type (
	// Config is the structure that holds the configuration variables of the application
	Config struct {
		Debug      bool   `envconfig:"DEBUG" default:"false"`
		SlackToken string `envconfig:"SLACK_TOKEN" default:""`
		DB         struct {
			DatabaseName string `envconfig:"POSTGRES_DB" default:""`
			Host         string `envconfig:"POSTGRES_HOST" default:""`
			User         string `envconfig:"POSTGRES_USER" default:""`
			Password     string `envconfig:"POSTGRES_PASSWORD" default:""`
		}
	}

	// EnvProcessFunc is the environment variable processor function
	EnvProcessFunc func(prefix string, spec interface{}) error
)

// LoadEnv loads config variables into Config struct
func LoadEnv(f EnvProcessFunc) (*Config, error) {
	var config Config
	err := f("", &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
