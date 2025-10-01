package main

import (
	"context"
	"fmt"
	"imageExtractor/internal/downloader"
	"testing"
	"time"
)

func TestDownloader(t *testing.T) {
	urls := []string{"https://app.akharinkhabar.ir/images/2025/09/30/5548fd5b-80af-406b-8d48-4355e4574447.jpeg"}
	d := downloader.NewHTTPDownloader(nil, 5)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	files, err := d.Download(ctx, urls, "./out")
	if err != nil {
		panic(err)
	}
	fmt.Println("downloaded:", files)
}
