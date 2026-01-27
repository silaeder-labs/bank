package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	gologger "github.com/nrf24l01/go-logger"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/silaeder-labs/bank/backend/postgres"
	"github.com/silaeder-labs/bank/backend/schemas"
)

func (h *Handler) CreateTransactionHandler(c echo.Context) error {
	req := c.Get("validatedBody").(*schemas.CreateTransactionRequest)
	from := c.Get("userID").(uuid.UUID)

	transaction := postgres.Transaction{
		From:        from,
		To:          req.TargetID,
		AmountCents: req.Amount,
		Description: req.Comment,
	}

	if err := transaction.Insert(h.DB, c.Request().Context()); err != nil {
		h.Logger.Log(gologger.LevelError, gologger.LogType("HTTP"), "Failed to create transaction: "+err.Error(), c.Get("traceId").(string))
		return c.JSON(500, echokitSchemas.GenError(c, echokitSchemas.INTERNAL_SERVER_ERROR, "failed to create transaction", nil))
	}

	return c.JSON(201, transaction.ToTransactionFull())
}