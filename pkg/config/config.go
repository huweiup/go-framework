package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load loads configuration from file and environment variables into the given struct.
func Load(path string, config interface{}) error {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		// Default search paths if no specific file is provided
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
