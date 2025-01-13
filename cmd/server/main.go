// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"retail-pulse/internal/api"
	"retail-pulse/internal/config"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
    cfg := config.Load()
    
    r := mux.NewRouter()
    api.SetupRoutes(r)
    
    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })
    
    handler := c.Handler(r)
    
    log.Printf("Server starting on port %s", cfg.Port)
    if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
        log.Fatal(err)
    }
}