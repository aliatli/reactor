package models

import "gorm.io/gorm"

type Edge struct {
	Source       string `json:"source"`
	Target       string `json:"target"`
	SourceHandle string `json:"sourceHandle"`
}

type State struct {
	gorm.Model
	Name               string           `gorm:"uniqueIndex"`
	PreliminaryActions []PrimitiveChain `gorm:"serializer:json"`
	MainAction         string
	PositionX          float64
	PositionY          float64
	SuccessTransition  string
	FailureTransition  string
	Edges              []Edge `gorm:"serializer:json"`
}

type PrimitiveChain struct {
	Primitives     []string
	ExecutionOrder int
}
