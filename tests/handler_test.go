package web

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"imageExtractor/internal/parser"
)

// mockExtractor is a mock implementation of the Extractor interface for testing
type mockExtractor struct {
	images  []string
	wantErr bool
}

func (m *mockExtractor) ExtractFromText(text string) ([]string, error) {
	if m.wantErr {
		return nil, fmt.Errorf("mock error")
	}
	return m.images, nil
}

func (m *mockExtractor) ExtractFromFile(path string) ([]string, error) {
	if m.wantErr {
		return nil, fmt.Errorf("mock error")
	}
	return m.images, nil
}

// TestHomeHandler tests the homeHandler function
func TestHomeHandler(t *testing.T) {
	// Load templates
	tmpl, err := template.ParseGlob("../../templates/*.html")
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   "این برنامه وب با استفاده از regex فایل‌های شما را استخراج می‌کند",
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := homeHandler(tmpl)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("homeHandler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			body := rr.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("homeHandler returned unexpected body: got %v want %v", body, tt.expectedBody)
			}
		})
	}
}

// TestExtractHandler tests the extractHandler function
func TestExtractHandler(t *testing.T) {
	// Load templates
	tmpl, err := template.ParseGlob("../../templates/*.html")
	if err != nil {
		t.Fatalf("Failed to parse templates: %v", err)
	}

	// Mock HTTP server for testing external requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<img src="/image.jpg"><img src="https://example.com/image.png">`))
	}))
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		formData       url.Values
		extractor      parser.Extractor
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid POST with URL",
			method:         http.MethodPost,
			formData:       url.Values{"link": {server.URL}},
			extractor:      &mockExtractor{images: []string{"/image.jpg", "https://example.com/image.png"}, wantErr: false},
			expectedStatus: http.StatusOK,
			expectedBody:   "2 تصویر یافت شد",
		},
		{
			name:           "POST with missing URL",
			method:         http.MethodPost,
			formData:       url.Values{},
			extractor:      &mockExtractor{},
			expectedStatus: http.StatusOK,
			expectedBody:   "URL is required",
		},
		{
			name:           "Invalid method (GET)",
			method:         http.MethodGet,
			formData:       url.Values{},
			extractor:      &mockExtractor{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
		{
			name:           "Invalid URL",
			method:         http.MethodPost,
			formData:       url.Values{"link": {"invalid://url"}},
			extractor:      &mockExtractor{},
			expectedStatus: http.StatusOK,
			expectedBody:   "Invalid URL format",
		},
		{
			name:           "Extractor error",
			method:         http.MethodPost,
			formData:       url.Values{"link": {server.URL}},
			extractor:      &mockExtractor{wantErr: true},
			expectedStatus: http.StatusOK,
			expectedBody:   "Failed to extract images",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/submit", strings.NewReader(tt.formData.Encode()))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := extractHandler(tmpl, tt.extractor)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("extractHandler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			body := rr.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("extractHandler returned unexpected body: got %v want %v", body, tt.expectedBody)
			}
		})
	}
}
