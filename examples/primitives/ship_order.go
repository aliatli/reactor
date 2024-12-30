package primitives

import (
	"fmt"

	"github.com/aliatli/reactor/internal/core"
)

type ShipOrder struct {
	// You might want to inject shipping service client here
	// shippingService ShippingService
}

func (s *ShipOrder) Execute(context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	// Verify all required data is present
	trackingNumber, ok := context.Data["trackingNumber"].(string)
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "missing tracking number"},
		}, nil
	}

	allocations, ok := context.Data["allocations"].(map[string]string)
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "missing inventory allocations"},
		}, nil
	}

	// Simulate shipping process
	// In a real implementation, you would:
	// 1. Notify warehouse
	// 2. Update inventory
	// 3. Notify shipping carrier
	// 4. Update order status

	shipmentID := fmt.Sprintf("SHIP-%s", trackingNumber)

	return &core.PrimitiveResult{
		Success: true,
		Data: map[string]interface{}{
			"shipmentID":     shipmentID,
			"shippingStatus": "in_transit",
			"shippedAt":      "2024-01-30T12:00:00Z", // Use actual time in production
			"allocations":    allocations,
		},
	}, nil
}
