package config

import (
	utilsConfig "github.com/nrf24l01/go-web-utils/config"
)

type Config struct {
	WebAppConfig   *utilsConfig.WebAppConfig
	PGConfig       *utilsConfig.PGConfig
}

func BuildConfigFromEnv() (*Config, error) {
	config := &Config{
		WebAppConfig:   utilsConfig.LoadWebAppConfigFromEnv(),
		PGConfig:       utilsConfig.LoadPGConfigFromEnv(),
	}

	return config, nil
}
