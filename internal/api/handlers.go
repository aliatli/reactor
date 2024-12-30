package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aliatli/reactor/internal/core"
)

func (s *Server) handleSaveFlow(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST /api/flow - Saving flow configuration")
	var flow struct {
		States map[string]core.StateDefinition `json:"states"`
	}

	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		log.Printf("Error decoding flow: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.stateDefinitions = flow.States
	log.Printf("Saved flow with %d states", len(flow.States))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (s *Server) handleGetStates(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /api/states - Returning %d states", len(s.stateDefinitions))
	json.NewEncoder(w).Encode(s.stateDefinitions)
}

func (s *Server) handleSaveState(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST /api/states - Saving new state")
	var state core.StateDefinition
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		log.Printf("Error decoding state: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Initialize all required fields
	if state.PreliminaryActions == nil {
		state.PreliminaryActions = []core.PrimitiveChain{}
	}
	if state.Position.X == 0 && state.Position.Y == 0 {
		state.Position = core.Position{
			X: float64(len(s.stateDefinitions)) * 100,
			Y: float64(len(s.stateDefinitions)) * 100,
		}
	}
	if state.Transitions.Success == "" {
		state.Transitions.Success = "none"
	}
	if state.Transitions.Failure == "" {
		state.Transitions.Failure = "none"
	}

	s.stateDefinitions[state.Name] = state
	log.Printf("Saved state: %s at position (%f, %f)", state.Name, state.Position.X, state.Position.Y)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"state":  state,
	})
}

func (s *Server) handleGetPrimitives(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /api/primitives - Returning available primitives")
	primitives := []string{
		"validateOrder",
		"checkInventory",
		"processPayment",
		"allocateInventory",
		"generateShippingLabel",
		"shipOrder",
	}
	json.NewEncoder(w).Encode(primitives)
}
