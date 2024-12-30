package api

import (
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
