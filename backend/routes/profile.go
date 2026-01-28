package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/middleware"
)

func RegisterProfileRoutes(e *echo.Group, h *handlers.Handler) {
	g := e.Group("/profile")
	g.Use(middleware.JWTMiddleware(h))
	g.GET("/me", h.GetBalanceHandler)
}