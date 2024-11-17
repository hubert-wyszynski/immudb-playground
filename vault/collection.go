package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://vault.immudb.io/ics/api/v1"

type CollectionStructure struct {
	Fields      []Field `json:"fields"`
	IdFieldName string  `json:"idFieldName"`
	Indexes     []Index `json:"indexes"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Index struct {
	Fields   []string `json:"fields"`
	IsUnique bool     `json:"isUnique"`
}

func EnsureCollection(apiKey string) error {
	url := fmt.Sprintf("%s/ledger/default/collection/default", baseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-API-Key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	collection := CollectionStructure{
		Fields: []Field{
			{Name: "account_number", Type: "STRING"},
			{Name: "account_name", Type: "STRING"},
			{Name: "iban", Type: "STRING"},
			{Name: "address", Type: "STRING"},
			{Name: "amount", Type: "DOUBLE"},
			{Name: "type", Type: "STRING"},
		},
		IdFieldName: "_id",
		Indexes: []Index{
			{
				Fields:   []string{"account_number"},
				IsUnique: true,
			},
		},
	}

	return createCollection(apiKey, collection)
}

func createCollection(apiKey string, collection CollectionStructure) error {
	url := fmt.Sprintf("%s/ledger/default/collection/default", baseURL)
	jsonData, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("error marshaling collection: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating collection: status code %d", resp.StatusCode)
	}

	return nil
}
