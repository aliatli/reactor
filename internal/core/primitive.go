package core

// PrimitiveResult represents the result of a primitive operation
type PrimitiveResult struct {
	Success   bool
	NextState string
	Data      map[string]interface{}
}

// Primitive defines the interface for all primitive operations
type Primitive interface {
	Execute(context *ExecutionContext) (*PrimitiveResult, error)
}
