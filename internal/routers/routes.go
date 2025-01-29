package routers

import (
	"journal-lite/internal/database"

	"github.com/labstack/echo/v4"
	"journal-lite/internal/auth"
	"journal-lite/internal/health"
)

func RegisterRoutes(e *echo.Echo) {

	database.Initialize()

	e.GET("/health", health.HealthCheckHandler(database.Db))

	e.POST("/login", auth.LoginHandler(database.Db))
	e.POST("/refresh-token", auth.RefreshTokenHandler())

	e.POST("/cookies/login", auth.LoginHandler(database.Db))
	e.POST("/cookies/refresh-token", auth.RefreshTokenHandler())
	e.POST("/cookies/logout", auth.LoginHandler(database.Db))
}
