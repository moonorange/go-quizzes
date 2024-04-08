package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

// APIWorker sends post requests
type APIWorker struct{}

func NewAPIWorker() *APIWorker {
	return &APIWorker{}
}

// Work waits for a few seconds and print a received URL.
func (p *APIWorker) Work(job *Job) {
	t := time.NewTimer(time.Duration(rand.Intn(5)) * time.Second)
	defer t.Stop()
	<-t.C

	server := mockServer()
	// Make a fake HTTP request to the mock server
	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf(`{"key": %s}`, job.Payload)))
	if err != nil {
		// Output error logs so that failed payload can be retried
		fmt.Printf("Error ocurred with payload %s: %+v", job.Payload, err)
		return
	}
	defer resp.Body.Close()
}

func mockServer() *httptest.Server {
	// Create a new HTTP test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Print the request body (fake processing)
		fmt.Println("Received POST request with body:", string(body))
	}))
	return server
}
