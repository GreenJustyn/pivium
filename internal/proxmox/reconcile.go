package proxmox

import (
	"fmt"
	"log"
	"os"

	"pivium/internal/config"
)

// Reconcile ensures the state of Proxmox resources matches the desired configuration.
func Reconcile(cfg config.Config) error {
	fmt.Println(">> Reconciling Proxmox State...")

	if !cfg.Proxmox.Enabled {
		fmt.Println("   Proxmox is disabled in the configuration.")
		return nil
	}

	// 1. Create a new Proxmox client.
	apiURL := os.Getenv("PROXMOX_API_URL")
	if apiURL == "" {
		return fmt.Errorf("PROXMOX_API_URL environment variable not set")
	}
	client, err := NewProxmoxClient(apiURL)
	if err != nil {
		return err
	}

	// 2. Get the current state of VMs and LXC containers from the Proxmox API.
	currentState, err := client.GetResources()
	if err != nil {
		return err
	}

	// 3. Get the desired state from the configuration.
	desiredState := cfg.Proxmox.Resources

	// 4. Compare the current state with the desired state.
	changes := compareStates(currentState, desiredState)
	log.Printf("Changes: %+v\n", changes)

	// 5. Apply the necessary changes (create, update, delete).
	if err := client.ApplyChanges(changes); err != nil {
		return err
	}

	fmt.Println("   Proxmox reconciliation complete.")
	return nil
}
