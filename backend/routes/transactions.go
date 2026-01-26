package routes

import (
	"github.com/labstack/echo/v4"
	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/middleware"
	"github.com/silaeder-labs/bank/backend/schemas"
)

func RegisterTransactionRoutes(e *echo.Group, h *handlers.Handler) {
	e.POST("/transactions", h.CreateTransactionHandler, echokitMw.BodyValidationMiddleware(func() interface{} {
		return &schemas.CreateTransactionRequest{}
	}), middleware.JWTMiddleware(h))
}