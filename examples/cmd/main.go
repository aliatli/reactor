package main

import (
	"fmt"
	"log"

	"github.com/aliatli/reactor/examples/primitives"
	"github.com/aliatli/reactor/internal/core"
	"github.com/aliatli/reactor/internal/executor"
)

func main() {
	// Create state executor
	stateExecutor := executor.NewStateExecutor()

	// Register all primitives
	primitives.RegisterPrimitives(stateExecutor.ChainExecutor.PrimitiveRegistry)

	// Define states
	orderReceivedState := core.StateDefinition{
		Name: "OrderReceived",
		PreliminaryActions: []core.PrimitiveChain{
			{
				Primitives:     []string{"validateOrder", "checkInventory"},
				ExecutionOrder: 1,
			},
		},
		MainAction: "processPayment",
	}
	orderReceivedState.Transitions.Success = "OrderFulfillment"
	orderReceivedState.Transitions.Failure = "OrderCancelled"

	orderFulfillmentState := core.StateDefinition{
		Name: "OrderFulfillment",
		PreliminaryActions: []core.PrimitiveChain{
			{
				Primitives:     []string{"allocateInventory"},
				ExecutionOrder: 1,
			},
			{
				Primitives:     []string{"generateShippingLabel"},
				ExecutionOrder: 2,
			},
		},
		MainAction: "shipOrder",
	}
	orderFulfillmentState.Transitions.Success = "OrderCompleted"
	orderFulfillmentState.Transitions.Failure = "CustomerServiceReview"

	// Register states
	stateExecutor.StateDefinitions["OrderReceived"] = orderReceivedState
	stateExecutor.StateDefinitions["OrderFulfillment"] = orderFulfillmentState

	// Create sample order data
	sampleOrder := map[string]interface{}{
		"id":     "ORD-12345",
		"amount": 99.99,
		"items": []interface{}{
			map[string]interface{}{
				"id":       "ITEM-1",
				"quantity": 1,
				"price":    49.99,
			},
			map[string]interface{}{
				"id":       "ITEM-2",
				"quantity": 1,
				"price":    50.00,
			},
		},
		"shippingAddress": map[string]interface{}{
			"street":  "123 Main St",
			"city":    "Springfield",
			"state":   "IL",
			"zipCode": "62701",
			"country": "USA",
		},
	}

	// Create execution context
	context := core.NewExecutionContext()
	context.Data["order"] = sampleOrder

	// Process the order through states
	currentState := "OrderReceived"
	fmt.Printf("Starting state machine with state: %s\n", currentState)

	for {
		fmt.Printf("\nExecuting state: %s\n", currentState)
		nextState, err := stateExecutor.ExecuteState(currentState, context)
		if err != nil {
			log.Fatalf("Error executing state %s: %v", currentState, err)
		}

		// Print relevant context data for debugging
		fmt.Printf("State execution completed. Context data:\n")
		printRelevantContextData(context)

		if nextState == "" {
			fmt.Printf("\nState machine completed at state: %s\n", currentState)
			break
		}

		if nextState == "OrderCancelled" || nextState == "CustomerServiceReview" {
			fmt.Printf("\nOrder requires attention: %s\n", nextState)
			break
		}

		if nextState == "OrderCompleted" {
			fmt.Printf("\nOrder successfully completed!\n")
			break
		}

		currentState = nextState
	}
}

func printRelevantContextData(context *core.ExecutionContext) {
	relevantKeys := []string{
		"orderValidated",
		"itemsAvailable",
		"paymentProcessed",
		"transactionID",
		"inventoryAllocated",
		"allocations",
		"trackingNumber",
		"labelURL",
		"shipmentID",
		"shippingStatus",
	}

	for _, key := range relevantKeys {
		if value, exists := context.Data[key]; exists {
			fmt.Printf("- %s: %v\n", key, value)
		}
	}
}
