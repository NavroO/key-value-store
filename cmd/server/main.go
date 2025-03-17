package main

import (
	"fmt"

	"github.com/NavroO/go-key-value-store/pkg/api"
	"github.com/NavroO/go-key-value-store/pkg/storage"
)

func main() {
	store := storage.NewInMemoryStorage()
	server := api.NewServer(store)

	port := "8080"
	fmt.Println("Starting server on port", port)
	server.Run(port)
}
