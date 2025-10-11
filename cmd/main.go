package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"imageExtractor/internal/parser"
	"imageExtractor/internal/web"
)

// main initializes the server and sets up routes
func main() {
	// Create a new HTMLExtractor instance
	extractor := parser.NewHTMLExtractor()

	// Parse HTML templates from the templates directory
	tmplPath := filepath.Join("templates", "*.html")
	templates, err := template.ParseGlob(tmplPath)
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// Set up HTTP multiplexer
	mux := http.NewServeMux()

	// Register handlers with dependencies
	web.RegisterHandlers(mux, templates, extractor)

	// Serve static files (e.g., CSS) from the templates directory
	fs := http.FileServer(http.Dir("templates"))
	mux.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// Start the server
	addr := ":8080"
	log.Printf("Server started at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
