package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // default
	}

	// 2. Init viper
	v := viper.New()

	// 3. Set config file dynamically
	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")

	// 4. Add search path
	v.AddConfigPath("./config") // relative to project root

	// 5. Enable ENV override
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 6. Read file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	// 7. Unmarshal
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return &cfg, nil
}
