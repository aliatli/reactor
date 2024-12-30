package api

import (
	"net/http"

	"github.com/aliatli/reactor/internal/core"
	"github.com/aliatli/reactor/internal/db"
	"github.com/gorilla/mux"
)

type Server struct {
	router           *mux.Router
	stateDefinitions map[string]core.StateDefinition
	db               *db.Database
}

func NewServer(database *db.Database) *Server {
	s := &Server{
		router:           mux.NewRouter(),
		stateDefinitions: make(map[string]core.StateDefinition),
		db:               database,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// Add CORS middleware
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	s.router.HandleFunc("/api/states", s.handleGetStates).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/states", s.handleSaveState).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/api/primitives", s.handleGetPrimitives).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/flow", s.handleSaveFlow).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/api/states/{name}", s.handleDeleteState).Methods("DELETE", "OPTIONS")
}

func (s *Server) Router() *mux.Router {
	return s.router
}
