package handlers

import (
	"github.com/lestrrat-go/jwx/v3/jwk"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/nrf24l01/go-web-utils/pgkit"
	"github.com/silaeder-labs/bank/backend/config"
)

type Handler struct {
	DB     *pgkit.DB
	Config *config.Config
	Jwks   *jwk.Cache
	Logger *gologger.Logger
}
