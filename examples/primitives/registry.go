package primitives

import (
	"github.com/aliatli/reactor/internal/core"
)

// RegisterPrimitives registers all available primitives in the chain executor
func RegisterPrimitives(registry map[string]core.Primitive) {
	registry["validateOrder"] = &ValidateOrder{}
	registry["checkInventory"] = &CheckInventory{}
	registry["processPayment"] = &ProcessPayment{}
	registry["allocateInventory"] = &AllocateInventory{}
	registry["generateShippingLabel"] = &GenerateShippingLabel{}
	registry["shipOrder"] = &ShipOrder{}
}
