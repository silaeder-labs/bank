package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	gologger "github.com/nrf24l01/go-logger"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/silaeder-labs/bank/backend/postgres"
	"github.com/silaeder-labs/bank/backend/schemas"
)

func (h *Handler) CreatePaymentHandler(c echo.Context) error {
	req := c.Get("validatedBody").(*schemas.CreatePaymentRequest)

	payment := postgres.Payment{
		From:        req.FromID,
		To:          req.ToID,
		Amount:      req.Amount,
		Description: req.Description,
		Status: 	string(schemas.StatusPending),
	}

	if err := payment.Insert(h.DB, c.Request().Context()); err != nil {
		h.Logger.Log(gologger.LevelError, gologger.LogType("DB"), "Failed to create payment: "+err.Error(), c.Get("traceId").(string))
		return c.JSON(http.StatusInternalServerError, echokitSchemas.GenError(c, echokitSchemas.INTERNAL_SERVER_ERROR, "failed to create payment", nil))
	}

	return c.JSON(http.StatusCreated, payment.ToPaymentFull())
}

func (h *Handler) GetPaymentHandler(c echo.Context) error {
	paymentIdStr := c.Param("uuid")
	userID := c.Get("userID").(uuid.UUID)
	
	paymentUUID, err := uuid.Parse(paymentIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echokitSchemas.GenError(c, echokitSchemas.BAD_REQUEST, "invalid payment ID", nil))
	}

	payment, err := postgres.GetPaymentByID(h.DB, c.Request().Context(), paymentUUID, userID)
	if err != nil {
		h.Logger.Log(gologger.LevelError, gologger.LogType("DB"), "Failed to get payment: "+err.Error(), c.Get("traceId").(string))
		return c.JSON(http.StatusInternalServerError, echokitSchemas.GenError(c, echokitSchemas.INTERNAL_SERVER_ERROR, "failed to get payment", nil))
	}

	return c.JSON(http.StatusOK, payment.ToPaymentFull())
}