package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"journal-lite/internal/database"
	"journal-lite/internal/posts"
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

	database.Initialize()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/feed", func(c echo.Context) error {

		posts, err := GetPostsByAccountID(1)

		if err != nil {
			posts = []Post{}
		}

		return c.Render(200, "feed-page", posts)
	})

	e.GET("/posts", func(c echo.Context) error {

		posts, err := GetPostsByAccountID(1)

		if err != nil {
			posts = []Post{}
		}

		return c.Render(200, "posts", posts)
	})

	e.GET("/open-delete-modal", func(c echo.Context) error {
		return c.Render(200, "delete-modal", nil)
	})

	e.GET("/open-create-modal", func(c echo.Context) error {
		return c.Render(200, "create-modal", nil)
	})

	e.GET("/close-modal", func(c echo.Context) error {
		return c.Render(200, "empty-div", nil)
	})

	e.POST("/create-post", func(c echo.Context) error {
		newPost := posts.Post{
			Content:   c.FormValue("content"),
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			AccountId: 1,
		}

		createdPost, err := posts.CreatePost(database.Db, newPost)
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "Error creating post")
		}

		return c.Render(201, "created-post-successfully", createdPost)
	})

	e.GET("/edit/:id", func(c echo.Context) error {
		return c.Render(200, "edit-post", nil)
	})

	e.GET("health", func(c echo.Context) error {

		response := map[string]string{
			"status": "healthy",
		}

		return c.JSON(200, response)
	})

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "user" && password == "pass" {
			token := "mock_jwt_token"

			cookie := http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour),
				Path:    "/",
			}

			c.SetCookie(&cookie)

			return c.Redirect(http.StatusFound, "/feed")
		}

		message := LoginBoxMessage{
			IsInvalidAttempt: true,
			Message:          "Invalid Credentials",
		}

		return c.Render(http.StatusUnauthorized, "index", message)
	})

	e.DELETE("/logout", func(c echo.Context) error {
		cookie := http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
			Path:    "/",
		}

		c.SetCookie(&cookie)

		return c.HTML(http.StatusOK, `<script>window.location.href = "/";</script>`)
	})

	e.Start(":8080")
}

type Post struct {
	Id        int    `db:"id"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	AccountId int    `db:"account_id"`
}

type LoginBoxMessage struct {
	IsInvalidAttempt bool
	Message          string
}

func GetPostsByAccountID(accountID int) ([]Post, error) {
	if database.Db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `
        SELECT id, content, created_at, updated_at, account_id
        FROM posts
        WHERE account_id = ?
				ORDER BY created_at DESC;
    `

	rows, err := database.Db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(
			&p.Id,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.AccountId,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return posts, nil
}
