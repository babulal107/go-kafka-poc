package main

import (
	"context"
	"errors"
	"github.com/babulal107/go-kafka-poc/producer/handler"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	g := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: g,
	}

	g.POST("/api/v2/order", handler.Order)

	// Run the server in a goroutine so it doesn’t block
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	GraceFulShutDown(srv)

}

func GraceFulShutDown(srv *http.Server) {

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Catch system signals

	<-quit // Wait for a shutdown signal
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited properly")
}
