# ImmuDB Playground API

A simple REST API service that manages account records using ImmuDB Vault.

## Running the Server

The server runs on port 8080 by default.
```bash
go run main.go
```

## Prerequisites

Before running the application, make sure you have:
- Set the `VAULT_API_KEY` environment variable with your ImmuDB Vault API key

## API Endpoints

### Add Account
Creates a new account record.

POST `/account`

**Request Body:**

```json
{
"account_number": "123456789",
"balance": 1000.50,
"currency": "USD",
"owner_name": "John Doe"
}
```

**Responses:**
- `201 Created`: Account successfully created
  ```json
  {
      "message": "Account record created successfully",
      "id": "document-id-from-vault"
  }
  ```
- `409 Conflict`: Account with given account number already exists
  ```json
  {
      "message": "Account with given account_number already exists"
  }
  ```
- `400 Bad Request`: Invalid request body
- `500 Internal Server Error`: Server-side error

### Get Accounts
Retrieves all account records.

GET `/accounts`

**Response:**
- `200 OK`: Returns a list of accounts from the vault
  ```json
  {
      "revisions": [
          {
              "documentId": "document-id",
              "revision": "revision-number",
              "document": {
                  "account_number": "123456789",
                  "balance": 1000.50,
                  "currency": "USD",
                  "owner_name": "John Doe"
              }
          }
          // ... more accounts
      ],
      "page": 1,
      "perPage": 100,
      "total": 1
  }
  ```
- `500 Internal Server Error`: Server-side error