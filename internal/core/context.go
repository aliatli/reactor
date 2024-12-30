package core

// ExecutionContext holds the shared state during execution
type ExecutionContext struct {
	Data map[string]interface{}
}

// NewExecutionContext creates a new execution context
func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		Data: make(map[string]interface{}),
	}
}
