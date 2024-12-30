package primitives

import (
	"fmt"

	"github.com/aliatli/reactor/internal/core"
)

type GenerateShippingLabel struct {
	// You might want to inject a shipping service client here
	// shippingService ShippingService
}

func (g *GenerateShippingLabel) Execute(context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	order, ok := context.Data["order"].(map[string]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid order format"},
		}, nil
	}

	// Extract shipping address
	shippingAddress, ok := order["shippingAddress"].(map[string]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid shipping address"},
		}, nil
	}

	// Simulate shipping label generation
	// In a real implementation, you would integrate with a shipping provider
	trackingNumber := fmt.Sprintf("TRACK-%s", order["id"])
	labelURL := fmt.Sprintf("https://shipping.example.com/labels/%s-%v.pdf",
		trackingNumber,
		shippingAddress,
	)

	return &core.PrimitiveResult{
		Success: true,
		Data: map[string]interface{}{
			"shippingLabelGenerated": true,
			"trackingNumber":         trackingNumber,
			"labelURL":               labelURL,
		},
	}, nil
}
