package parser

import (
	"io/ioutil"
	"regexp"
)

// Extractor defines the interface for text and file-based image extraction.
type Extractor interface {
	ExtractFromText(text string) ([]string, error)
	ExtractFromFile(path string) ([]string, error)
}

// HTMLExtractor implements the Extractor interface for HTML content.
type HTMLExtractor struct{}

// NewHTMLExtractor returns a new instance of HTMLExtractor.
func NewHTMLExtractor() *HTMLExtractor {
	return &HTMLExtractor{}
}

// ExtractFromText extracts image URLs from the given HTML text.
func (h *HTMLExtractor) ExtractFromText(text string) ([]string, error) {
	re, err := regexp.Compile(`([^"]*\.(jpg|png|jpeg|gif))`)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllString(text, -1)
	return matches, nil
}

// ExtractFromFile reads an HTML file and extracts image URLs.
func (h *HTMLExtractor) ExtractFromFile(path string) ([]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return h.ExtractFromText(string(data))
}
