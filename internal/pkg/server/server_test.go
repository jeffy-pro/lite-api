package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	for _, port := range []string{"8081", ":8081"} {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Create a mock handler
		handler := http.NewServeMux()
		handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		go ServeHTTP(ctx, port, handler) // Use a different port to avoid conflicts

		time.Sleep(100 * time.Millisecond) // Give the server some time to start

		// Test if the server is up and running
		resp, err := http.Get("http://localhost:8081")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		// Cancel the context to simulate shutdown
		cancel()

		time.Sleep(100 * time.Millisecond) // Give the server some time to shut down

		// Test if the server has shut down
		resp, err = http.Get("http://localhost:8081")
		assert.Error(t, err)
	}

}
