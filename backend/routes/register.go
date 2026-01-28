package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/silaeder-labs/bank/backend/handlers"
)

func RegisterRoutes(e *echo.Group, h *handlers.Handler) {
	RegisterTransactionRoutes(e, h)
	RegisterProfileRoutes(e, h)
}
