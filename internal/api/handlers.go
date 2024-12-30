package api

import (
	"encoding/json"
	"net/http"

	"github.com/aliatli/reactor/internal/core"
	"github.com/gorilla/mux"
)

type Server struct {
	router           *mux.Router
	stateDefinitions map[string]core.StateDefinition
}

func NewServer() *Server {
	s := &Server{
		router:           mux.NewRouter(),
		stateDefinitions: make(map[string]core.StateDefinition),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/api/states", s.handleGetStates).Methods("GET")
	s.router.HandleFunc("/api/states", s.handleSaveState).Methods("POST")
	s.router.HandleFunc("/api/primitives", s.handleGetPrimitives).Methods("GET")
	s.router.HandleFunc("/api/flow", s.handleSaveFlow).Methods("POST")
}

func (s *Server) handleSaveFlow(w http.ResponseWriter, r *http.Request) {
	var flow struct {
		States map[string]core.StateDefinition `json:"states"`
	}

	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save flow configuration
	s.stateDefinitions = flow.States

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
