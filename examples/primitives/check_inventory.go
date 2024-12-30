package primitives

import (
	"github.com/aliatli/reactor/internal/core"
)

type CheckInventory struct {
	// You might want to inject a database client or inventory service here
	// inventoryService InventoryService
}

func (c *CheckInventory) Execute(context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	order, ok := context.Data["order"].(map[string]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid order format"},
		}, nil
	}

	// Simulate inventory check
	// In a real implementation, you would check against your inventory system
	items, ok := order["items"].([]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid items format"},
		}, nil
	}

	// Check each item's availability
	availabilityMap := make(map[string]bool)
	allAvailable := true

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemID, ok := itemMap["id"].(string)
		if !ok {
			continue
		}

		// Simulate inventory check - replace with actual inventory check
		isAvailable := true // This would be the result of your inventory query
		availabilityMap[itemID] = isAvailable
		allAvailable = allAvailable && isAvailable
	}

	return &core.PrimitiveResult{
		Success: allAvailable,
		Data: map[string]interface{}{
			"inventoryChecked": true,
			"itemsAvailable":   availabilityMap,
		},
	}, nil
}
