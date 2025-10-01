package parser

import (
	"io/ioutil"
	"regexp"
)

// Extractor interface
type Extractor interface {
	ExtractFromText(text string) ([]string, error)
	ExtractFromFile(path string) ([]string, error)
}

// HTMLExtractor concrete struct
type HTMLExtractor struct{}

// NewHTMLExtractor constructor
func NewHTMLExtractor() *HTMLExtractor {
	return &HTMLExtractor{}
}

// ExtractFromText extracts image links from HTML string
func (h *HTMLExtractor) ExtractFromText(text string) ([]string, error) {
	re, err := regexp.Compile(`([^"]*\.(jpg|png|jpeg|gif))`)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllString(text, -1)
	return matches, nil
}

// ExtractFromFile reads file and extracts image links
func (h *HTMLExtractor) ExtractFromFile(path string) ([]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return h.ExtractFromText(string(data))
}
