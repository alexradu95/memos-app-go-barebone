package routers

import (
	"journal-lite/internal/database"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"journal-lite/internal/accounts"
	"journal-lite/internal/auth"
	"journal-lite/internal/health"
	"journal-lite/internal/posts"
)

func RegisterRoutes(e *echo.Echo) {

	e.GET("/health", health.HealthCheckHandler(database.Db))

	e.POST("/login", auth.LoginHandler(database.Db))
	e.POST("/refresh-token", auth.RefreshTokenHandler())

	e.POST("/cookies/login", auth.LoginHandler(database.Db))
	e.POST("/cookies/refresh-token", auth.RefreshTokenHandler())
	e.POST("/cookies/logout", auth.LoginHandler(database.Db))

	e.POST("/accounts", accounts.CreateAccountHandler(database.Db))

	// protected routes
	protected := e.Group("/", middleware.AuthMiddleware)
	protected.DELETE("/accounts/:id", accounts.DeleteAccountHandler(database.Db))
	protected.GET("/posts", posts.GetPostsHandler(database.Db))
	protected.POST("/posts", posts.CreatePostHandler(database.Db))
	protected.DELETE("/posts/:id", posts.DeletePostHandler(database.Db))
	protected.PUT("/posts/:id", posts.UpdatePostHandler(database.Db))

}
