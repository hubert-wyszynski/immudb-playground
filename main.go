package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hubert-wyszynski/immudb-playground/handlers"
	"github.com/hubert-wyszynski/immudb-playground/vault"
)

func main() {
	apiKey := os.Getenv("VAULT_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set VAULT_API_KEY environment variable")
	}

	if err := vault.EnsureCollection(apiKey); err != nil {
		log.Fatalf("Failed to ensure collection exists: %v", err)
	}

	accountHandler := handlers.NewAccountHandler(apiKey)

	http.HandleFunc("/account", accountHandler.AddAccount)
	http.HandleFunc("/accounts", accountHandler.GetAccounts)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
