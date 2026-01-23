package handlers

import (
	"github.com/MicahParks/keyfunc"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/silaeder-labs/bank/backend/config"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Config *config.Config
	Jwks   *keyfunc.JWKS
	Logger *gologger.Logger
}
