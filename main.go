package main

import (
	"database/sql"
	"html/template"
	"journal-lite/internal/database"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
)

func main() {
	err := database.Initialize()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/load", LoadContent)

	http.ListenAndServe(":8080", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func LoadContent(w http.ResponseWriter, r *http.Request) {
	html := `<div class="p-4 bg-blue-100">
               <p>Content loaded via HTMX!</p>
             </div>`
	w.Write([]byte(html))
}
