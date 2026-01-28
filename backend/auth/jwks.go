package auth

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/silaeder-labs/bank/backend/config"
)

func RegisterJwks(cfg *config.KeyCloakConfig, logger *gologger.Logger, ctx *context.Context) (*jwk.Cache, error) {
	c, err := jwk.NewCache(*ctx, httprc.NewClient())
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("AUTH"), fmt.Sprintf("failed to create cache: %s", err), "")
		return nil, err
	}

	if err := c.Register(*ctx, cfg.URL); err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("AUTH"), fmt.Sprintf("failed to register google JWKS: %s", err), "")
		return nil, err
	}
	return c, nil
}
