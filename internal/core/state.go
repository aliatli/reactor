package core

// NextState represents the name of the next state
type NextState string

// StateDefinition defines the structure of a state
type StateDefinition struct {
	Name               string
	PreliminaryActions []PrimitiveChain
	MainAction         string
	Transitions        struct {
		Success NextState
		Failure NextState
	}
}

// PrimitiveChain represents a chain of primitive operations
type PrimitiveChain struct {
	Primitives     []string
	ExecutionOrder int
}
