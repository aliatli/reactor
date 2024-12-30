package api

import (
	"encoding/json"
	"net/http"

	"github.com/aliatli/reactor/internal/core"
)

func (s *Server) handleSaveFlow(w http.ResponseWriter, r *http.Request) {
	var flow struct {
		States map[string]core.StateDefinition `json:"states"`
	}

	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.stateDefinitions = flow.States

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (s *Server) handleGetStates(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.stateDefinitions)
}

func (s *Server) handleSaveState(w http.ResponseWriter, r *http.Request) {
	var state core.StateDefinition
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.stateDefinitions[state.Name] = state
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetPrimitives(w http.ResponseWriter, r *http.Request) {
	// Return list of available primitives
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
