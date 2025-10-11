package parser

import (
	"os"
	"reflect"
	"testing"
)

// TestNewHTMLExtractor tests the creation of a new HTMLExtractor
func TestNewHTMLExtractor(t *testing.T) {
	extractor := NewHTMLExtractor()
	if extractor == nil {
		t.Error("NewHTMLExtractor returned nil")
	}
	if _, ok := extractor.(*HTMLExtractor); !ok {
		t.Error("NewHTMLExtractor did not return an HTMLExtractor")
	}
}

// TestExtractFromText tests the ExtractFromText method
func TestExtractFromText(t *testing.T) {
	extractor := NewHTMLExtractor()
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name: "Valid HTML with images",
			input: `<img src="https://example.com/image1.jpg">
                    <img src="/upload/image2.png">
                    <a href="document.pdf">Link</a>
                    <img src="relative/image3.jpeg">`,
			expected: []string{
				"https://example.com/image1.jpg",
				"/upload/image2.png",
				"relative/image3.jpeg",
			},
			wantErr: false,
		},
		{
			name:     "Empty HTML",
			input:    `<div>No images here</div>`,
			expected: []string{},
			wantErr:  false,
		},
		{
			name:     "No matching extensions",
			input:    `<img src="https://example.com/file.pdf">`,
			expected: []string{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractor.ExtractFromText(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFromText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractFromText() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestExtractFromFile tests the ExtractFromFile method
func TestExtractFromFile(t *testing.T) {
	extractor := NewHTMLExtractor()

	// Create a temporary HTML file for testing
	tempFile, err := os.CreateTemp("", "test_*.html")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test content to the file
	content := `<img src="https://example.com/test.jpg"><img src="/images/test.png">`
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		path     string
		expected []string
		wantErr  bool
	}{
		{
			name: "Valid HTML file",
			path: tempFile.Name(),
			expected: []string{
				"https://example.com/test.jpg",
				"/images/test.png",
			},
			wantErr: false,
		},
		{
			name:     "Non-existent file",
			path:     "nonexistent.html",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractor.ExtractFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractFromFile() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
