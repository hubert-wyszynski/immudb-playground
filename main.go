package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hubert-wyszynski/immudb-playground/handlers"
	"github.com/hubert-wyszynski/immudb-playground/vault"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next(w, r)
    }
}

func main() {
	apiKey := os.Getenv("VAULT_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set VAULT_API_KEY environment variable")
	}

	if err := vault.EnsureCollection(apiKey); err != nil {
		log.Fatalf("Failed to ensure collection exists: %v", err)
	}

	accountHandler := handlers.NewAccountHandler(apiKey)

	http.HandleFunc("/account", enableCORS(accountHandler.AddAccount))
    http.HandleFunc("/accounts", enableCORS(accountHandler.GetAccounts))

    fmt.Println("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
