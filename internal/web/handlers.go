package web

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	"imageExtractor/internal/parser"
)

// RegisterHandlers sets up the HTTP routes
func RegisterHandlers(mux *http.ServeMux, templates *template.Template, extractor parser.Extractor) {
	mux.HandleFunc("/", homeHandler(templates))
	mux.HandleFunc("/submit", extractHandler(templates, extractor))
}

// homeHandler renders the homepage
func homeHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Render the index.html template
		if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

// extractHandler processes the provided URL and extracts images
func extractHandler(templates *template.Template, extractor parser.Extractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			renderError(templates, w, "Failed to parse form data")
			return
		}

		// Get the URL from the form
		link := r.FormValue("link")
		if link == "" {
			renderError(templates, w, "URL is required")
			return
		}

		// Add https:// if no protocol is specified
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			link = "https://" + link
		}

		// Validate and parse the URL
		parsedURL, err := url.ParseRequestURI(link)
		if err != nil {
			renderError(templates, w, "Invalid URL format")
			return
		}

		// Fetch the content from the URL
		resp, err := http.Get(parsedURL.String())
		if err != nil {
			renderError(templates, w, "Failed to fetch URL content")
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			renderError(templates, w, "Failed to read URL content")
			return
		}

		// Extract image URLs
		images, err := extractor.ExtractFromText(string(body))
		if err != nil {
			renderError(templates, w, "Failed to extract images")
			return
		}

		// Convert relative URLs to absolute URLs
		absoluteImages := make([]string, 0, len(images))
		for _, img := range images {
			absoluteURL := resolveRelativeURL(parsedURL, img)
			absoluteImages = append(absoluteImages, absoluteURL)
		}

		// Prepare data for the result template
		data := struct {
			Images []struct {
				Link string
			}
		}{
			Images: make([]struct{ Link string }, len(absoluteImages)),
		}
		for i, img := range absoluteImages {
			data.Images[i].Link = img
		}

		// Render the result template
		if err := templates.ExecuteTemplate(w, "result.html", data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

// resolveRelativeURL converts a relative URL to an absolute URL based on the base URL
func resolveRelativeURL(baseURL *url.URL, relative string) string {
	// If the URL is already absolute (starts with http:// or https://), return it
	if strings.HasPrefix(relative, "http://") || strings.HasPrefix(relative, "https://") {
		return relative
	}

	// Parse the relative URL
	parsedRelative, err := url.Parse(relative)
	if err != nil {
		return relative // Return as-is if parsing fails
	}

	// Resolve the relative URL against the base URL
	absoluteURL := baseURL.ResolveReference(parsedRelative)
	return absoluteURL.String()
}

// renderError renders the error template with the given message
func renderError(templates *template.Template, w http.ResponseWriter, message string) {
	data := struct {
		Error string
	}{
		Error: message,
	}
	if err := templates.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
