package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"journal-lite/internal/accounts"
	"journal-lite/internal/auth"
	"journal-lite/internal/database"
	"journal-lite/internal/posts"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, r *http.Request) error {
	return t.templates.ExecuteTemplate(w, name+".html", data)
}

func newTemplate() *Template {
	funcMap := template.FuncMap{
		"formatDate": func(date string) string {
			t, err := time.Parse(time.RFC3339, date)
			if err != nil {
				return date
			}
			return t.Format("January 2, 2006")
		},
	}
	// Use Must to panic on template parsing errors
	return &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
	}
}

var templates = newTemplate()

func main() {
	database.Initialize()
	defer database.CloseDB() // Ensure DB closes when main exits.

	http.HandleFunc("/", handler)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Middleware-like behavior for all requests:
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	// Simple routing
	switch r.URL.Path {
	case "/":
		if r.Method == http.MethodGet {
			renderTemplate(w, r, "index", nil)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/feed":
		if r.Method == http.MethodGet {
			authMiddleware(http.HandlerFunc(feedHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/posts":
		if r.Method == http.MethodGet {
			authMiddleware(http.HandlerFunc(postsHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/register":
		switch r.Method {
		case http.MethodGet:
			renderTemplate(w, r, "register-box", nil)
		case http.MethodPost:
			registerHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/open-delete-modal":
		authMiddleware(http.HandlerFunc(openDeleteModalHandler)).ServeHTTP(w, r)

	case "/open-edit-modal":
		authMiddleware(http.HandlerFunc(openEditModalHandler)).ServeHTTP(w, r)

	case "/posts/update":
		if r.Method == http.MethodPatch {
			authMiddleware(http.HandlerFunc(updatePostHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/posts/delete":
		if r.Method == http.MethodDelete {
			authMiddleware(http.HandlerFunc(deletePostHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/open-create-modal":
		if r.Method == http.MethodGet {
			authMiddleware(http.HandlerFunc(openCreateModalHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/close-modal":
		if r.Method == http.MethodGet {
			renderTemplate(w, r, "empty-div", nil)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/create-post":
		if r.Method == http.MethodPost {
			authMiddleware(http.HandlerFunc(createPostHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/search":
		if r.Method == http.MethodGet {
			authMiddleware(http.HandlerFunc(searchHandler)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/health":
		if r.Method == http.MethodGet {
			healthHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/login":
		if r.Method == http.MethodPost {
			loginHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	case "/logout":
		if r.Method == http.MethodDelete {
			logoutHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	default:
		http.NotFound(w, r)
	}
}

// --- Handlers ---

func feedHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := GetPostsByAccountID(1) //  Get the account ID from the context.
	if err != nil {
		handleError(w, r, "Error fetching posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, r, "feed-page", posts)
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := GetPostsByAccountID(1) //  Get the account ID from the context.
	if err != nil {
		handleError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, r, "posts", posts)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")

	if password != passwordConfirmation {
		message := LoginBoxMessage{
			IsInvalidAttempt: true,
			Message:          "Passwords do not match",
		}
		renderTemplate(w, r, "register-box", message)
		return
	}

	newAccount := accounts.Account{
		Username:     username,
		PasswordHash: password,
	}

	_, err := accounts.CreateAccount(database.Db, newAccount)
	if err != nil {
		message := LoginBoxMessage{
			IsInvalidAttempt: true,
			Message:          "Failed to create account: " + err.Error(),
		}
		renderTemplate(w, r, "register-box", message)
		return
	}

	message := LoginBoxMessage{
		IsInvalidAttempt: false,
		Message:          "Account created successfully",
	}

	renderTemplate(w, r, "register-account-complete", message)
}

func openDeleteModalHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		handleError(w, r, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleError(w, r, "Invalid ID format", http.StatusBadRequest)
		return
	}

	post := posts.GetPost(database.Db, 1, id) // Replace 1 with actual user ID
	renderTemplate(w, r, "delete-modal", post)
}

func openEditModalHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		handleError(w, r, "ID required.", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleError(w, r, "Id Invalid.", http.StatusBadRequest)
		return
	}

	post := posts.GetPost(database.Db, 1, id) // Get the account ID from the context
	renderTemplate(w, r, "edit-modal", post)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Could not parse the form.", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		handleError(w, r, "ID required.", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleError(w, r, "Invalid ID format", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	err = posts.UpdatePost(database.Db, content, id)
	if err != nil {
		handleError(w, r, "Could not update post.", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, r, "empty-div", nil)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		handleError(w, r, "ID required.", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleError(w, r, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = posts.DeletePost(database.Db, id)
	if err != nil {
		handleError(w, r, "Error deleting post.", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, r, "empty-div", nil)
}

func openCreateModalHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "create-modal", nil)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Could not create the post.", http.StatusBadRequest)
		return
	}
	newPost := posts.Post{
		Content:   r.FormValue("content"),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		AccountId: 1, //  Get the account ID from the context.
	}

	createdPost, err := posts.CreatePost(database.Db, newPost)
	if err != nil {
		http.Error(w, "Error creating post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, r, "created-post-successfully", createdPost)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	params := posts.QueryParams{
		AccountId:  1,                           // Get the account ID from context
		SearchText: r.URL.Query().Get("search"), // This is how you get query params in net/http
		PageSize:   10,
		PageNumber: 1,
	}

	posts := posts.GetPosts(database.Db, params)
	renderTemplate(w, r, "feed", posts)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	err := database.Db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status": "error", "database": "error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "ok", "database": "ok"}`))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, "Could not parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := auth.Login(database.Db, username, password)
	if err != nil {
		message := LoginBoxMessage{
			IsInvalidAttempt: true,
			Message:          err.Error(),
		}
		renderTemplate(w, r, "index", message) //Render login page with error
		return
	}

	if token != "" {
		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			Secure:   true, //  Send over HTTPS only
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/feed", http.StatusFound)
		return
	}

	message := LoginBoxMessage{
		IsInvalidAttempt: true,
		Message:          "Invalid Credentials",
	}

	renderTemplate(w, r, "index", message) // Render login page
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expired time
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/") // Use HX-Redirect for HTMX
	w.WriteHeader(http.StatusOK)       // Send 200 OK
}

// --- Helper Functions ---

func renderTemplate(w http.ResponseWriter, r *http.Request, tmplName string, data interface{}) {
	err := templates.Render(w, tmplName, data, r)
	if err != nil {
		// Log the error properly.  Don't just print to stdout.
		log.Printf("Error rendering template %s: %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	log.Println(message) // Log the error
	w.WriteHeader(statusCode)
	// Consider rendering a simple error template here, if appropriate.
	if r.Header.Get("HX-Request") == "true" {
		// For HTMX, render a simple error message.
		_, _ = w.Write([]byte(fmt.Sprintf(`<div class="error">%s</div>`, message)))
	} else {
		// For regular requests, return a simple text error.
		http.Error(w, message, statusCode)
	}
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

// authMiddleware is a middleware function to protect routes.
func authMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Redirect to login for HTMX and standard requests.
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/")
					return
				}
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			handleError(w, r, "Error reading cookie", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims, err := auth.ValidateToken(tokenStr) // Use the new function
		if err != nil {
			// Clear the invalid cookie
			clearCookie := http.Cookie{
				Name:    "token",
				Value:   "",
				Expires: time.Unix(0, 0),
				Path:    "/",
			}
			http.SetCookie(w, &clearCookie)

			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/") // Redirect using HTMX header
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		//  Store user ID in context (example)
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx)) // Pass the context to the next handler
	})
}
