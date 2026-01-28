package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type KeyCloakConfig struct {
	Realm      string `env:"KEYCLOAK_REALM"`
	AuthServer string `env:"KEYCLOAK_AUTH_SERVER"`
	ISSUER_URL string `env:"KEYCLOAK_ISSUER_URL" envDefault:""`
	URL        string
}

func LoadKeyCloakConfigFromEnv() *KeyCloakConfig {
	config := &KeyCloakConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	if config.ISSUER_URL == "" {
		config.ISSUER_URL = fmt.Sprintf("%s/realms/%s", config.AuthServer, config.Realm)
	}
	config.URL = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", config.AuthServer, config.Realm)
	return config
}
