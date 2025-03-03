package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"astralis/internal/adapters/primary/rest"
	"astralis/internal/adapters/secondary/astronomyapi"
	"astralis/internal/adapters/secondary/nasaapi"
	"astralis/internal/core/ports"
	"astralis/internal/core/service"
	"astralis/pkg/config"
)

func main() {
	l := log.New(os.Stdout, "[Astralis API] ", 3)

	c := config.LoadConfig()
	l.Printf("loading config...")

	var repositories []ports.EventRepository
	astronomyRepo := astronomyapi.NewAstronomyAPIRepository()
	repositories = append(repositories, astronomyRepo)
	l.Printf("loading AstronomyAPI...")

	if c.NasaAPIKey() != "" {
		nasaRepo := nasaapi.NewNASARepository(c.NasaAPIKey())
		repositories = append(repositories, nasaRepo)
		l.Printf("loading NasaAPI...")
	}

	// Initialize service
	eventService := service.NewEventService(repositories)

	// Initialize REST handler
	handler := rest.NewHandler(eventService)

	// Create router and register routes
	router := gin.Default()
	handler.RegisterRoutes(router)

	l.Printf("Server starting on port %s", c.APIPort())
	srv := &http.Server{
		Addr:    c.APIPort(),
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("server crashed: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		l.Println("timeout of 5 seconds.")
	}
	l.Println("Server exiting")
}
