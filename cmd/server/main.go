package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NavroO/go-key-value-store/pkg/api"
	"github.com/NavroO/go-key-value-store/pkg/storage"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using default configuration")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	store := storage.NewInMemoryStorage()
	server := api.NewServer(store)

	_, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Println("Starting server on port 8080")
		server.Run(":8080")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Println("\nReceived signal:", sig)

	cancel()
	store.Shutdown()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Shutdown error:", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
