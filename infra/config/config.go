package config

// cspell:ignore mapstructure

type Config struct {
	App struct {
		Name string `yaml:"name" mapstructure:"name"`
		Port int    `yaml:"port" mapstructure:"port"`
	}
	Database struct {
		DbURL string `yaml:"db_url" mapstructure:"db_url"`
	}
}
