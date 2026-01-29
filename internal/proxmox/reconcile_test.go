package proxmox_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"pivium/internal/config"
	"pivium/internal/proxmox"
)

func TestReconcile(t *testing.T) {
	// Mock Proxmox API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	// Set the PROXMOX_API_URL environment variable to the mock server's URL
	os.Setenv("PROXMOX_API_URL", server.URL)

	// Create a sample config
	cfg := config.Config{
		Proxmox: struct {
			Enabled   bool                       `json:"enabled"`
			Role      string                     `json:"role"`
			Resources []config.ProxmoxResource `json:"resources"`
		}{
			Enabled: true,
			Resources: []config.ProxmoxResource{
				{
					Type: "qemu",
					VMID: 100,
				},
			},
		},
	}

	// Call the Reconcile function
	if err := proxmox.Reconcile(cfg); err != nil {
		t.Errorf("Reconcile failed: %v", err)
	}
}
