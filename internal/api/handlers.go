package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aliatli/reactor/internal/core"
	"github.com/aliatli/reactor/internal/models"
	"github.com/gorilla/mux"
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

	// Save each state to the database
	for _, stateDefinition := range flow.States {
		// Convert primitive chains
		modelChains := make([]models.PrimitiveChain, len(stateDefinition.PreliminaryActions))
		for i, chain := range stateDefinition.PreliminaryActions {
			modelChains[i] = models.PrimitiveChain{
				Primitives:     chain.Primitives,
				ExecutionOrder: chain.ExecutionOrder,
			}
		}

		// Convert edges
		modelEdges := make([]models.Edge, len(stateDefinition.Edges))
		for i, edge := range stateDefinition.Edges {
			modelEdges[i] = models.Edge{
				Source:       edge.Source,
				Target:       edge.Target,
				SourceHandle: edge.SourceHandle,
			}
		}

		state := &models.State{
			Name:               stateDefinition.Name,
			PreliminaryActions: modelChains,
			MainAction:         stateDefinition.MainAction,
			PositionX:          stateDefinition.Position.X,
			PositionY:          stateDefinition.Position.Y,
			SuccessTransition:  stateDefinition.Transitions.Success,
			FailureTransition:  stateDefinition.Transitions.Failure,
			Edges:              modelEdges,
		}

		if err := s.db.SaveState(state); err != nil {
			log.Printf("Error saving state %s: %v", state.Name, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.stateDefinitions = flow.States
	log.Printf("Saved flow with %d states", len(flow.States))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func (s *Server) handleGetStates(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /api/states - Fetching states from database")

	states, err := s.db.GetAllStates()
	if err != nil {
		log.Printf("Error fetching states: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var stateDefinitions []core.StateDefinition
	for _, state := range states {
		// Convert PrimitiveChain slice
		primitiveChains := make([]core.PrimitiveChain, len(state.PreliminaryActions))
		for i, chain := range state.PreliminaryActions {
			primitiveChains[i] = core.PrimitiveChain{
				Primitives:     chain.Primitives,
				ExecutionOrder: chain.ExecutionOrder,
			}
		}

		// Convert Edge slice
		edges := make([]core.Edge, len(state.Edges))
		for i, edge := range state.Edges {
			edges[i] = core.Edge{
				Source:       edge.Source,
				Target:       edge.Target,
				SourceHandle: edge.SourceHandle,
			}
		}

		stateDef := core.StateDefinition{
			Name:               state.Name,
			PreliminaryActions: primitiveChains,
			Position: core.Position{
				X: state.PositionX,
				Y: state.PositionY,
			},
			Edges: edges,
			Transitions: struct {
				Success string `json:"success"`
				Failure string `json:"failure"`
			}{
				Success: state.SuccessTransition,
				Failure: state.FailureTransition,
			},
		}
		stateDefinitions = append(stateDefinitions, stateDef)
	}

	log.Printf("Returning states: %v", stateDefinitions)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stateDefinitions)
}

func (s *Server) handleSaveState(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST /api/states - Saving new state")
	var stateDefinition core.StateDefinition
	if err := json.NewDecoder(r.Body).Decode(&stateDefinition); err != nil {
		log.Printf("Error decoding state: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to database model
	modelChains := make([]models.PrimitiveChain, len(stateDefinition.PreliminaryActions))
	for i, chain := range stateDefinition.PreliminaryActions {
		modelChains[i] = models.PrimitiveChain{
			Primitives:     chain.Primitives,
			ExecutionOrder: chain.ExecutionOrder,
		}
	}

	// Convert edges
	modelEdges := make([]models.Edge, len(stateDefinition.Edges))
	for i, edge := range stateDefinition.Edges {
		modelEdges[i] = models.Edge{
			Source:       edge.Source,
			Target:       edge.Target,
			SourceHandle: edge.SourceHandle,
		}
	}

	state := &models.State{
		Name:               stateDefinition.Name,
		PreliminaryActions: modelChains,
		MainAction:         stateDefinition.MainAction,
		PositionX:          stateDefinition.Position.X,
		PositionY:          stateDefinition.Position.Y,
		SuccessTransition:  stateDefinition.Transitions.Success,
		FailureTransition:  stateDefinition.Transitions.Failure,
		Edges:              modelEdges,
	}

	if err := s.db.SaveState(state); err != nil {
		log.Printf("Error saving state: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.stateDefinitions[state.Name] = stateDefinition
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"state":  stateDefinition,
	})
}

func (s *Server) handleDeleteState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stateName := vars["name"]
	log.Printf("DELETE /api/states/%s - Deleting state", stateName)

	if err := s.db.DeleteState(stateName); err != nil {
		log.Printf("Error deleting state: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(s.stateDefinitions, stateName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
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
