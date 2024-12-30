package primitives

import (
	"github.com/aliatli/reactor/internal/core"
)

type ValidateOrder struct{}

func (v *ValidateOrder) Execute(context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	order, ok := context.Data["order"].(map[string]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid order format"},
		}, nil
	}

	// Add basic validation logic
	if _, exists := order["id"]; !exists {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "order id is required"},
		}, nil
	}

	return &core.PrimitiveResult{
		Success: true,
		Data:    map[string]interface{}{"orderValidated": true},
	}, nil
}
