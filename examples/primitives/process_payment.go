package primitives

import (
	"github.com/aliatli/reactor/internal/core"
)

type ProcessPayment struct {
	// You might want to inject a payment service client here
	// paymentService PaymentService
}

func (p *ProcessPayment) Execute(context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	order, ok := context.Data["order"].(map[string]interface{})
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid order format"},
		}, nil
	}

	amount, ok := order["amount"].(float64)
	if !ok {
		return &core.PrimitiveResult{
			Success: false,
			Data:    map[string]interface{}{"error": "invalid amount"},
		}, nil
	}

	// Simulate payment processing
	// In a real implementation, you would integrate with a payment gateway
	paymentSuccessful := amount > 0 // Replace with actual payment processing

	result := &core.PrimitiveResult{
		Success: paymentSuccessful,
		Data: map[string]interface{}{
			"paymentProcessed": true,
			"transactionID":    "txn_123", // This would be from your payment provider
			"amount":           amount,
		},
	}

	if !paymentSuccessful {
		result.Data["error"] = "payment failed"
	}

	return result, nil
}
