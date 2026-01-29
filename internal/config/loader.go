package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProxmoxResource represents a single VM or LXC container.
type ProxmoxResource struct {
	Type     string `json:"type"`
	VMID     int    `json:"vmid"`
	Hostname string `json:"hostname"`
	Template string `json:"template"`
	Memory   int    `json:"memory"`
	Cores    int    `json:"cores"`
	Net0     string `json:"net0"`
	Storage  string `json:"storage"`
	State    string `json:"state"`
}

// Config represents the aggregated state of the node
type Config struct {
	Version      string `json:"version"`
	ClusterName  string `json:"cluster_name"`
	UpdatePolicy string `json:"update_policy"`
	System       struct {
		Timezone string   `json:"timezone"`
		Packages []string `json:"packages"`
	} `json:"system"`
	Proxmox struct {
		Enabled   bool              `json:"enabled"`
		Role      string            `json:"role"`
		Resources []ProxmoxResource `json:"resources"`
	} `json:"proxmox"`
	Ceph struct {
		Enabled bool   `json:"enabled"`
		Device  string `json:"device"`
	} `json:"ceph"`
}

// Load cascades JSON files to build the final config object
func Load(rootDir string, hostname string) (*Config, error) {
	cfg := &Config{}

	// Order of precedence: Lowest -> Highest
	// 1. Defaults
	// 2. Host Specific
	files := []string{
		filepath.Join(rootDir, "configs", "defaults.json"),
		filepath.Join(rootDir, "configs", "hosts", hostname+".json"),
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// Skip if host file doesn't exist, but defaults must exist
			if file == files[0] {
				return nil, fmt.Errorf("defaults.json missing at %s", file)
			}
			continue
		}

		fmt.Printf("Loading config layer: %s\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		// JSON Unmarshal writes over existing fields in the struct
		if err := json.Unmarshal(content, cfg); err != nil {
			return nil, fmt.Errorf("malformed JSON in %s: %w", file, err)
		}
	}

	return cfg, nil
}