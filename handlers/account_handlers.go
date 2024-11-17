package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hubert-wyszynski/immudb-playground/models"
)

const (
	baseURL           = "https://vault.immudb.io/ics/api/v1"
	defaultLedger     = "default"
	defaultCollection = "default"
	defaultPerPage    = 100
)

type AccountHandler struct {
	apiKey string
	client *http.Client
}

func NewAccountHandler(apiKey string) *AccountHandler {
	return &AccountHandler{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (h *AccountHandler) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := h.buildURL(endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func (h *AccountHandler) buildURL(endpoint string) string {
	return fmt.Sprintf("%s/ledger/%s/collection/%s/%s",
		baseURL, defaultLedger, defaultCollection, endpoint)
}

func (h *AccountHandler) AddAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var record models.AccountRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.makeRequest(r.Context(), http.MethodPut, "document", record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Account with given account_number already exists",
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to store record", resp.StatusCode)
		return
	}

	var vaultResp models.VaultResponse
	if err := json.NewDecoder(resp.Body).Decode(&vaultResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account record created successfully",
		"id":      vaultResp.DocumentID,
	})
}

func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	searchReq := map[string]interface{}{
		"page":    1,
		"perPage": defaultPerPage,
	}

	resp, err := h.makeRequest(r.Context(), http.MethodPost, "documents/search", searchReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
