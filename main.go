package main

import (
	"net/http"
)

type Entry struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	DateCreated string `json:"date_created"`
}

func main() {
	http.HandleFunc("/", HomePage)
	http.ListenAndServe(":8080", nil)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
