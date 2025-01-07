package main

import (
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name+".html", data)
}

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	posts := []Post{
		{Date: "2020-01-01", Content: "Hello, World!"},
		{Date: "2020-01-02", Content: "Hello, Echo!"},
		{Date: "2020-01-01", Content: "Hello, World!"},
		{Date: "2020-01-02", Content: "Hello, Echo!"},
		{Date: "2020-01-01", Content: "Hello, World!"},
		{Date: "2020-01-02", Content: "Hello, Echo!"},
		{Date: "2020-01-01", Content: "Hello, World!"},
		{Date: "2020-01-01", Content: "Hello, World!"},
		{Date: "2020-01-02", Content: "Hello, Echo!"},
		{Date: "2020-01-02", Content: "Hello, Echo!"},
		{Date: "2020-01-01", Content: "Hello, World!"},
		{Date: "2020-01-02", Content: "Hello, Echo!"},
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/feed", func(c echo.Context) error {
		return c.Render(200, "feed-page", posts)
	})

	e.GET("/posts", func(c echo.Context) error {
		return c.Render(200, "posts", posts)
	})

	e.GET("/create-post", func(c echo.Context) error {
		return c.Render(200, "create-post", nil)
	})

	e.GET("/edit/:id", func(c echo.Context) error {
		return c.Render(200, "edit-post", nil)
	})

	e.POST("/login", func(c echo.Context) error {
		// Authenticate user (validate credentials)
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "user" && password == "pass" {
			// Generate JWT or session token
			token := "mock_jwt_token"

			// Set cookie with JWT
			cookie := http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour),
				Path:    "/",
			}

			c.SetCookie(&cookie)

			// Redirect or render feed block
			return c.Redirect(http.StatusFound, "/feed")
		}

		// Render login error block with htmx
		return c.Render(http.StatusUnauthorized, "login-box", map[string]interface{}{
			"Error": "Invalid credentials",
		})
	})

	e.DELETE("/logout", func(c echo.Context) error {
		// Clear cookie
		cookie := http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
			Path:    "/",
		}

		c.SetCookie(&cookie)

		// Redirect or render login block
		return c.HTML(http.StatusOK, `<script>window.location.href = "/";</script>`)
	})

	e.Start(":8080")
}

type Post struct {
	Date    string
	Content string
}
