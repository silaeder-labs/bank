package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		return c.JSON(http.StatusInternalServerError, echokitSchemas.GenError(c, echokitSchemas.INTERNAL_SERVER_ERROR, "failed to create transaction", nil))
	}

	return c.JSON(http.StatusCreated, transaction.ToTransactionFull())
}

func (h *Handler) GetTransactionsHandler(c echo.Context) error {
	req := c.Get("validatedQuery").(*schemas.GetTransactionsRequest)
	userID := c.Get("userID").(uuid.UUID)

	offset := (req.Page - 1) * req.Size

	transactions, err := postgres.GetTransactionsByUserID(h.DB, c.Request().Context(), userID, req.Size, offset)
	if err != nil {
		h.Logger.Log(gologger.LevelError, gologger.LogType("HTTP"), "Failed to get transactions: "+err.Error(), c.Get("traceId").(string))
		return c.JSON(http.StatusInternalServerError, echokitSchemas.GenError(c, echokitSchemas.INTERNAL_SERVER_ERROR, "failed to get transactions", nil))
	}
	
	var resp []schemas.TransactionFull
	for _, t := range transactions {
		resp = append(resp, t.ToTransactionFull())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTransactionByIDHandler(c echo.Context) error {
	userID := c.Get("userID").(uuid.UUID)
	transactionIDStr := c.Param("uuid")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echokitSchemas.GenError(c, echokitSchemas.BAD_REQUEST, "invalid transaction ID", nil))
	}

	transaction, err := postgres.GetTransactionByID(h.DB, c.Request().Context(), transactionID, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echokitSchemas.GenError(c, echokitSchemas.NOT_FOUND, "transaction not found", nil))
		}
		h.Logger.Log(gologger.LevelError, gologger.LogType("HTTP"), "Failed to get transaction: "+err.Error(), c.Get("traceId").(string))
		return c.JSON(http.StatusInternalServerError, echokitSchemas.GenError(c, echokitSchemas.INTERNAL_SERVER_ERROR, "failed to get transaction", nil))
	}

	if transaction.From != userID && transaction.To != userID {
		return c.JSON(http.StatusForbidden, echokitSchemas.GenError(c, echokitSchemas.FORBIDDEN, "access to transaction denied", nil))
	}

	return c.JSON(http.StatusOK, transaction.ToTransactionFull())
}