package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIsPortOpen(t *testing.T) {
	nu := NetworkIOUtils{}

	t.Run("Port is open", func(t *testing.T) {
		// Start a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer server.Close()

		address := server.Listener.Addr().String()
		if !nu.IsPortOpen(address) {
			t.Errorf("Expected port %s to be open", address)
		}
	})

	t.Run("Port is closed", func(t *testing.T) {
		address := "localhost:12345"
		if nu.IsPortOpen(address) {
			t.Errorf("Expected port %s to be closed", address)
		}
	})
}

func TestDownloadNetworkConfiguration(t *testing.T) {
	nu := NetworkIOUtils{}
	algodDataDir := os.TempDir()

	t.Run("Successful download", func(t *testing.T) {
		// Mock http.Get
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"genesis": "data"}`))
			if err != nil {
				return
			}
		}))
		defer server.Close()

		err := nu.DownloadNetworkConfiguration(server.URL, algodDataDir)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Failed download", func(t *testing.T) {
		// Mock http.Get to return a non-200 status code
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		err := nu.DownloadNetworkConfiguration(server.URL, algodDataDir)
		if err == nil {
			t.Errorf("Expected error, got none")
		}
	})
}
