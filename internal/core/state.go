package core

// NextState represents the name of the next state
type NextState string

// StateDefinition defines the structure of a state
type StateDefinition struct {
	Name               string           `json:"name"`
	PreliminaryActions []PrimitiveChain `json:"preliminaryActions"`
	MainAction         string           `json:"mainAction,omitempty"`
	Position           Position         `json:"position"`
	Transitions        struct {
		Success string `json:"success"`
		Failure string `json:"failure"`
	} `json:"transitions"`
}

// PrimitiveChain represents a chain of primitive operations
type PrimitiveChain struct {
	Primitives     []string
	ExecutionOrder int
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
