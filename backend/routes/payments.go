package routes

import (
	"github.com/labstack/echo/v4"
	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/middleware"
	"github.com/silaeder-labs/bank/backend/schemas"
)

func RegisterPaymentsRoutes(e *echo.Group, h *handlers.Handler) {
	g := e.Group("/payments")
	g.POST("", h.CreatePaymentHandler, middleware.JWTMiddleware(h, true), echokitMw.BodyValidationMiddleware(func() interface{} {
		return &schemas.CreatePaymentRequest{}
	}))
}