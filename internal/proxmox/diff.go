package proxmox

import "pivium/internal/config"

// Changes holds the lists of resources to be created, updated, and deleted.
type Changes struct {
	Create []config.ProxmoxResource
	Update []config.ProxmoxResource
	Delete []config.ProxmoxResource
}

// compareStates compares the current state with the desired state and returns the changes.
func compareStates(current, desired []config.ProxmoxResource) Changes {
	changes := Changes{}
	currentMap := make(map[int]config.ProxmoxResource)

	for _, resource := range current {
		currentMap[resource.VMID] = resource
	}

	for _, desiredResource := range desired {
		if _, ok := currentMap[desiredResource.VMID]; !ok {
			changes.Create = append(changes.Create, desiredResource)
		} else {
			// Resource exists, check for updates
			// For now, we'll just add it to the update list if it exists.
			// In the future, we can add more sophisticated diffing logic here.
			changes.Update = append(changes.Update, desiredResource)
			delete(currentMap, desiredResource.VMID)
		}
	}

	// Any remaining resources in currentMap are to be deleted.
	for _, resource := range currentMap {
		changes.Delete = append(changes.Delete, resource)
	}

	return changes
}
