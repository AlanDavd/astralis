package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"astralis/internal/adapters/primary/rest"
	"astralis/internal/adapters/secondary/astronomyapi"
	"astralis/internal/adapters/secondary/nasaapi"
	"astralis/internal/core/ports"
	"astralis/internal/core/service"
)

func main() {
	// Initialize repositories
	nasaRepo := nasaapi.NewNASARepository(
		os.Getenv("NASA_API_KEY"),
	)

	astronomyRepo := astronomyapi.NewAstronomyAPIRepository()

	// Create a slice of repositories
	var repositories []ports.EventRepository
	repositories = append(repositories, nasaRepo, astronomyRepo)

	// Initialize service
	eventService := service.NewEventService(repositories)

	// Initialize REST handler
	handler := rest.NewHandler(eventService)

	// Create router and register routes
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Using repositories: NASA API, Visible Planets API")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
} 