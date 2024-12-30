package executor

import (
	"fmt"

	"github.com/aliatli/reactor/internal/core"
)

type PrimitiveChainExecutor struct {
	PrimitiveRegistry map[string]core.Primitive
}

func NewPrimitiveChainExecutor() *PrimitiveChainExecutor {
	return &PrimitiveChainExecutor{
		PrimitiveRegistry: make(map[string]core.Primitive),
	}
}

func (pce *PrimitiveChainExecutor) Execute(chain core.PrimitiveChain, context *core.ExecutionContext) (*core.PrimitiveResult, error) {
	for _, primitiveName := range chain.Primitives {
		primitive, exists := pce.PrimitiveRegistry[primitiveName]
		if !exists {
			return nil, fmt.Errorf("primitive not found: %s", primitiveName)
		}

		result, err := primitive.Execute(context)
		if err != nil {
			return nil, err
		}

		if !result.Success {
			return result, nil // Break chain on failure
		}

		// Update context with result data
		for k, v := range result.Data {
			context.Data[k] = v
		}
	}

	return &core.PrimitiveResult{Success: true}, nil
}
