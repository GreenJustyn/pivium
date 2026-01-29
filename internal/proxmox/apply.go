package proxmox

import (
	"fmt"

	"pivium/internal/config"
)

// ApplyChanges applies the changes to the Proxmox resources.
func (c *ProxmoxClient) ApplyChanges(changes Changes) error {
	for _, resource := range changes.Create {
		if err := c.createResource(resource); err != nil {
			return err
		}
	}

	for _, resource := range changes.Update {
		if err := c.updateResource(resource); err != nil {
			return err
		}
	}

	for _, resource := range changes.Delete {
		if err := c.deleteResource(resource); err != nil {
			return err
		}
	}

	return nil
}

func (c *ProxmoxClient) createResource(resource config.ProxmoxResource) error {
	fmt.Printf("Creating resource: %+v\n", resource)
	// TODO: Implement resource creation logic.
	return nil
}

func (c *ProxmoxClient) updateResource(resource config.ProxmoxResource) error {
	fmt.Printf("Updating resource: %+v\n", resource)
	// TODO: Implement resource update logic.
	return nil
}

func (c *ProxmoxClient) deleteResource(resource config.ProxmoxResource) error {
	fmt.Printf("Deleting resource: %+v\n", resource)
	// TODO: Implement resource deletion logic.
	return nil
}
