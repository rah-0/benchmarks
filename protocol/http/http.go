package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	httpAddress = "localhost:8080"
	client      *Client
)

func httpServerStart() {
	http.HandleFunc("/test", handleHTTPRequest)

	log.Printf("HTTP server listening on %s", httpAddress)
	if err := http.ListenAndServe(httpAddress, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		log.Printf("Error reading request body: %v", err)
		return
	}
	defer r.Body.Close()

	// Process the received bytes (optional)
	_ = body // Do nothing with the data

	w.WriteHeader(http.StatusOK)
}

type Client struct {
	client  *http.Client
	url     string
	reader  *bytes.Reader
	headers http.Header
}

func NewHTTPClient() error {
	client = &Client{
		client: &http.Client{},
		url:    "http://" + httpAddress + "/test",
		reader: bytes.NewReader(nil),
		headers: http.Header{
			"Content-Type": []string{"application/octet-stream"},
		},
	}

	return nil
}

func (c *Client) sendBytes(data []byte) error {
	c.reader.Reset(data)

	req, err := http.NewRequest(http.MethodPost, c.url, c.reader)
	if err != nil {
		return err
	}

	req.Header = c.headers
	req.ContentLength = int64(len(data))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	return nil
}
