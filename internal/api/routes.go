// internal/api/routes.go
package api

import (
	"retail-pulse/internal/api/handlers"
	"retail-pulse/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
    // Apply common middleware
    r.Use(middleware.LoggingMiddleware)
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/submit", handlers.SubmitJob).Methods("POST")
    api.HandleFunc("/status", handlers.GetJobStatus).Methods("GET")
}