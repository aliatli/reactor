package primitives

import (
	"github.com/aliatli/reactor/internal/core"
)

type AllocateInventory struct {
	// You might want to inject an inventory service client here
	// inventoryService InventoryService
}

func (a *AllocateInventory) Execute(context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	order, ok := context.Data["order"].(map[string]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid order format"},
		}, nil
	}

	// Check if inventory was previously verified
	itemsAvailable, ok := context.Data["itemsAvailable"].(map[string]bool)
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "inventory check not performed"},
		}, nil
	}

	// Simulate inventory allocation
	// In a real implementation, you would reserve inventory in your system
	allAllocated := true
	allocatedItems := make(map[string]string)

	items, ok := order["items"].([]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid items format"},
		}, nil
	}

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemID, ok := itemMap["id"].(string)
		if !ok {
			continue
		}

		if available, exists := itemsAvailable[itemID]; exists && available {
			// Simulate allocation - replace with actual inventory allocation
			allocatedItems[itemID] = "alloc_" + itemID
		} else {
			allAllocated = false
			break
		}
	}

	return &core.PrimitiveResult{
		Success: allAllocated,
		Data: map[string]interface{}{
			"inventoryAllocated": allAllocated,
			"allocations":        allocatedItems,
		},
	}, nil
}
