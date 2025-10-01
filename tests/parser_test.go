package main

import (
	"fmt"
	"imageExtractor/internal/parser"
	"testing"
)

func TestParser(t *testing.T) {
	filePath := "test.html"
	extractor := parser.NewHTMLExtractor()
	files, err := extractor.ExtractFromFile(filePath)
	if err != nil {
		fmt.Printf("ExtractFromFile error: %v\n", err)
		return
	}
	fmt.Println(len(files))
}
