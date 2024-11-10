package authguard

import (
	"net/http"
	"strings"

	common "github.com/abialemuel/billing-engine/pkg/common/http"

	"github.com/abialemuel/billing-engine/config"
	"github.com/labstack/echo/v4"
)

// Constants for handling JWT and keys
var (
	PrefixHeader      = "Bearer "
	PrefixHeaderBasic = "Basic "
)

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthGuard holds dependencies like API key and configuration
type AuthGuard struct {
	cfg config.MainConfig
}

// NewAuthGuard creates a new instance of AuthGuard
func NewAuthGuard(cfg config.MainConfig) *AuthGuard {
	return &AuthGuard{
		cfg: cfg,
	}
}

// Basic middleware validates Basic Auth tokens
func (g *AuthGuard) Basic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, PrefixHeaderBasic) {
			return c.JSON(http.StatusUnauthorized, common.NewUnauthorizedResponse("Authorization header missing/invalid"))
		}

		username, _, ok := c.Request().BasicAuth()
		if !ok {
			return c.JSON(http.StatusUnauthorized, common.NewUnauthorizedResponse("Invalid token"))
		}

		// set the service name to the context
		c.Set("username", username)

		return next(c)
	}
}
