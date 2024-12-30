package executor

import "github.com/aliatli/reactor/internal/core"

type StateExecutor struct {
	StateDefinitions map[string]core.StateDefinition
	ChainExecutor    *PrimitiveChainExecutor
}

func NewStateExecutor() *StateExecutor {
	return &StateExecutor{
		StateDefinitions: make(map[string]core.StateDefinition),
		ChainExecutor:    NewPrimitiveChainExecutor(),
	}
}

func (se *StateExecutor) ExecuteState(stateName string, context *core.ExecutionContext) (string, error) {
	state, exists := se.StateDefinitions[stateName]
	if !exists {
		return "", nil
	}

	// Execute preliminary actions in order
	for _, chain := range state.PreliminaryActions {
		result, err := se.ChainExecutor.Execute(chain, context)
		if err != nil {
			return string(state.Transitions.Failure), err
		}
		if !result.Success {
			return string(state.Transitions.Failure), nil
		}
	}

	// Execute main action if present
	if state.MainAction != "" {
		result, err := se.ChainExecutor.Execute(core.PrimitiveChain{
			Primitives: []string{state.MainAction},
		}, context)
		if err != nil {
			return string(state.Transitions.Failure), err
		}
		if !result.Success {
			return string(state.Transitions.Failure), nil
		}
	}

	return string(state.Transitions.Success), nil
}
