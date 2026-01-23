package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/middleware"
)

func RegisterRoutes(e *echo.Group, h *handlers.Handler) {
	e.GET("/check_auth", h.CheckAuth, middleware.JWTMiddleware(h))
}
