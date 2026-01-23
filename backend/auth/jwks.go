package auth

import (
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/silaeder-labs/bank/backend/config"
)

func RegisterJwks(cfg *config.KeyCloakConfig, logger *gologger.Logger) (*keyfunc.JWKS, error) {
	jwks, err := keyfunc.Get(cfg.URL, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			logger.Log(gologger.LevelFatal, gologger.LogType("AUTH"), fmt.Sprintf("ERROR IN JWKS REFRESH: %s", err), "")
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	})
	if err != nil {
		return nil, err
	}
	return jwks, nil
}
