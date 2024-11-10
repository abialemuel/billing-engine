package http

import (
	"github.com/abialemuel/billing-engine/pkg/common/http/middleware/authguard"

	"github.com/labstack/echo/v4"
)

// RegisterPath Register V1 API path
func RegisterPath(e *echo.Echo, h *Handler, authGuard *authguard.AuthGuard) {
	if h == nil {
		panic("item controller cannot be nil")
	}

	// Auth implementation
	e.GET("/v1/billing-engine/user/outstanding", h.GetOutstandingLoan, authGuard.Basic)
	e.POST("/v1/billing-engine/user/payment", h.MakePayment, authGuard.Basic)

}
