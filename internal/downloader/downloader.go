package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Downloader defines the contract for downloading multiple URLs.
type Downloader interface {
	Download(ctx context.Context, urls []string, path string) ([]string, error)
}

// httpDownloader implements Downloader using HTTP client and worker pool.
type httpDownloader struct {
	client     *http.Client
	maxWorkers int
}

// NewHTTPDownloader returns a new instance of httpDownloader with sane defaults.
func NewHTTPDownloader(client *http.Client, maxWorkers int) Downloader {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second, // default timeout
		}
	}
	if maxWorkers <= 0 {
		maxWorkers = 5 // default number of workers
	}
	if maxWorkers > 20 {
		maxWorkers = 20 // prevent overload
	}

	return &httpDownloader{
		client:     client,
		maxWorkers: maxWorkers,
	}
}

// Download concurrently downloads all URLs into the given directory.
// It returns a list of file paths for successfully downloaded files.
func (d *httpDownloader) Download(ctx context.Context, urls []string, path string) ([]string, error) {
	// Ensure the download directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("error creating download directory: %w", err)
	}

	var wg sync.WaitGroup
	urlChan := make(chan string, len(urls))    // queue of URLs
	resultChan := make(chan string, len(urls)) // file paths of successful downloads
	errChan := make(chan error, len(urls))     // errors during download

	// Worker goroutines
	for i := 0; i < d.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range urlChan {
				// Stop early if context is canceled
				select {
				case <-ctx.Done():
					return
				default:
				}

				// Download a single file
				filePath, err := d.downloadFile(ctx, u, path)
				if err != nil {
					// Non-blocking send: if ctx is done, discard error
					select {
					case errChan <- fmt.Errorf("failed to download %s: %w", u, err):
					case <-ctx.Done():
					}
					continue
				}

				// Non-blocking send result
				select {
				case resultChan <- filePath:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Goroutine to feed URLs into urlChan
	go func() {
		defer close(urlChan)
		for _, u := range urls {
			select {
			case <-ctx.Done():
				return
			case urlChan <- u:
			}
		}
	}()

	// Goroutine to close resultChan when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []string

	// Collect results until resultChan is closed or context is canceled
	for {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		case r, ok := <-resultChan:
			if !ok {
				// At this point all workers are finished
				// If any errors exist, return the first one (or aggregate if needed)
				select {
				case err := <-errChan:
					return results, err
				default:
					return results, nil
				}
			}
			results = append(results, r)
		}
	}
}

// downloadFile downloads a single URL into the target directory.
// It returns the local file path where the file was saved.
func (d *httpDownloader) downloadFile(ctx context.Context, url, dir string) (string, error) {
	// Build HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Perform HTTP request
	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, url)
	}

	// Derive filename from URL base
	filename := filepath.Base(url)
	// Fallback if filename is empty (e.g., URL ends with /)
	if filename == "" || filename == "/" {
		filename = fmt.Sprintf("file_%d", time.Now().UnixNano())
	}
	filePath := filepath.Join(dir, filename)

	// Create file on disk
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy response body into the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}
