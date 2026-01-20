package handlers

import (
	"github.com/silaeder-labs/bank/backend/config"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Config *config.Config
}
