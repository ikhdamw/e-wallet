package model

import "time"

type Transfer struct {
	ID            string    `json:"id" db:"id"`
	FromWalletID  string    `json:"from_wallet_id" db:"from_wallet_id"`
	ToWalletID    string    `json:"to_wallet_id" db:"to_wallet_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Status        string    `json:"status" db:"status"`
	Description   string    `json:"description" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type ExternalTransfer struct {
	ID                string    `json:"id" db:"id"`
	TransactionID     string    `json:"transaction_id" db:"transaction_id"`
	Provider          string    `json:"provider" db:"provider"`
	RecipientAccount  string    `json:"recipient_account" db:"recipient_account"`
	RecipientName     string    `json:"recipient_name" db:"recipient_name"`
	Amount            float64   `json:"amount" db:"amount"`
	Fee               float64   `json:"fee" db:"fee"`
	Status            string    `json:"status" db:"status"`
	ProviderReference string    `json:"provider_reference" db:"provider_reference"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

type InternalTransferRequest struct {
	RecipientEmail string  `json:"recipient_email" binding:"required,email"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Description    string  `json:"description"`
}

type ExternalTransferRequest struct {
	Provider         string  `json:"provider" binding:"required"`
	RecipientAccount string  `json:"recipient_account" binding:"required"`
	Amount           float64 `json:"amount" binding:"required,gt=0"`
	Description      string  `json:"description"`
}

type TransferResponse struct {
	TransferID string  `json:"transfer_id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	Fee        float64 `json:"fee,omitempty"`
}
