package amqp

import (
	"github.com/abialemuel/billing-engine/pkg/user/business"
)

type Handler struct {
	service business.UserService
}

func NewHandler(s business.UserService) *Handler {
	return &Handler{service: s}
}
