package models

type AccountRecord struct {
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
	IBAN          string  `json:"iban"`
	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
}

type VaultResponse struct {
	DocumentID    string `json:"documentId"`
	TransactionID string `json:"transactionId"`
}
