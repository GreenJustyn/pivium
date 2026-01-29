package ceph_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"pivium/internal/ceph"
	"pivium/internal/config"
)

func TestReconcile(t *testing.T) {
	// Mock Ceph API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	// Set the CEPH_API_URL environment variable to the mock server's URL
	os.Setenv("CEPH_API_URL", server.URL)

	// Create a sample config
	cfg := config.Config{
		Ceph: struct {
			Enabled bool   `json:"enabled"`
			Device  string `json:"device"`
		}{
			Enabled: true,
		},
	}

	// Call the Reconcile function
	if err := ceph.Reconcile(cfg); err != nil {
		t.Errorf("Reconcile failed: %v", err)
	}
}
