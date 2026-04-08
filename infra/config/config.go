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
	Redis struct {
		Host     string `yaml:"host" mapstructure:"host"`
		Port     int    `yaml:"port" mapstructure:"port"`
		Password string `yaml:"password" mapstructure:"password"`
		DB       int    `yaml:"db" mapstructure:"db"`
	}
}
