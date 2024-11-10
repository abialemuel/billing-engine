package http

import (
	"net/http"

	"github.com/abialemuel/billing-engine/config"
	common "github.com/abialemuel/billing-engine/pkg/common/http"
	"github.com/abialemuel/billing-engine/pkg/common/http/validator"
	"github.com/abialemuel/billing-engine/pkg/user/api/http/request"
	"github.com/abialemuel/billing-engine/pkg/user/api/http/response"
	"github.com/abialemuel/billing-engine/pkg/user/business"
	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service business.UserService
	config  *config.MainConfig
}

// NewHandler Construct user API handler
func NewHandler(service business.UserService, cfg *config.MainConfig) *Handler {
	return &Handler{
		service,
		cfg,
	}
}

// GetOutstandingLoan Get outstanding loan by user ID
func (h *Handler) GetOutstandingLoan(c echo.Context) error {
	ctx, span := apm.StartTransaction(c.Request().Context(), "Handler::GetOutstandingLoan")
	defer apm.EndTransaction(span)

	username := c.Get("username").(string)

	res, err := h.service.GetOutstandingLoan(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewValidationErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, response.NewOutstandingLoanInfoResponse(res))
}

// MakePayment for loan
func (h *Handler) MakePayment(c echo.Context) error {
	ctx, span := apm.StartTransaction(c.Request().Context(), "Handler::MakePayment")
	defer apm.EndTransaction(span)

	req := new(request.MakePaymentReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse("Invalid Body"))
	}
	if msg, check := validator.Validation(req); !check {
		return c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(msg))
	}

	username := c.Get("username").(string)

	err := h.service.MakePayment(ctx, username, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewValidationErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, common.NewDefaultSuccessResponse())
}
