package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get bearer token from cookie
		cookie, err := c.Cookie("bearerToken")
		if err != nil {
			// For HTMX requests, redirect to login page
			if c.Request().Header.Get("HX-Request") == "true" {
				return c.Redirect(http.StatusFound, "/")
			}
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		// Extract and validate claims
		claims, err := ExtractClaims(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusFound, "/")
		}

		if err := ValidateClaims(claims); err != nil {
			return c.Redirect(http.StatusFound, "/")
		}

		// Store claims in context for later use
		c.Set("claims", claims)

		// Continue to the next middleware/handler
		return next(c)
	}
}
