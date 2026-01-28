package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	gologger "github.com/nrf24l01/go-logger"
	"github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/silaeder-labs/bank/backend/postgres"
)

func (h *Handler) GetBalanceHandler(c echo.Context) error {
	uid := c.Get("userID").(uuid.UUID)
	
	balance, err := postgres.GetBalanceByUserID(h.DB, c.Request().Context(), uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, schemas.GenError(c, schemas.NOT_FOUND, "balance not found", nil))
		}
		h.Logger.Log(gologger.LevelError, gologger.LogType("DB"), "Failed to get balance: "+err.Error(), c.Get("traceId").(string))
		return c.JSON(http.StatusInternalServerError, schemas.GenError(c, schemas.INTERNAL_SERVER_ERROR, "Failed to get balance", nil))
	}

	balanceFull := balance.ToBalanceFull()

	return c.JSON(http.StatusOK, balanceFull)
}