package handlers

import (
	"github.com/lestrrat-go/jwx/v3/jwk"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/silaeder-labs/bank/backend/config"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Config *config.Config
	Jwks   *jwk.Cache
	Logger *gologger.Logger
}
