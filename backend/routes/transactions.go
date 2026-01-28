package routes

import (
	"github.com/labstack/echo/v4"
	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/middleware"
	"github.com/silaeder-labs/bank/backend/schemas"
)

func RegisterTransactionRoutes(e *echo.Group, h *handlers.Handler) {
	g := e.Group("/transactions")
	g.Use(middleware.JWTMiddleware(h, false))
	g.POST("", h.CreateTransactionHandler, echokitMw.BodyValidationMiddleware(func() interface{} {
		return &schemas.CreateTransactionRequest{}
	}))
	g.GET("", h.GetTransactionsHandler, echokitMw.QueryValidationMiddleware(func() interface{} {
		return &schemas.GetTransactionsRequest{}
	}))
	g.GET("/:uuid", h.GetTransactionByIDHandler, echokitMw.PathUuidV4Middleware("uuid"))
}
