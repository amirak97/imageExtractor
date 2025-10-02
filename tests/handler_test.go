package main

import (
	"imageExtractor/internal/web"
	"net/http"
	"testing"
)

func Test_handler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", web.HomeHandler)
	mux.HandleFunc("/submit", web.SubmitHandler)

	// استفاده از middleware
	handler := web.SessionMiddleware(mux)

	http.ListenAndServe(":8080", handler)
}
