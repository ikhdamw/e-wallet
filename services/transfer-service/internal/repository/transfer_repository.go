package repository

import (
	"database/sql"
	"time"

	"github.com/ikhdamw/e-wallet/transfer-service/internal/model"
	"github.com/google/uuid"
)

type TransferRepository interface {
	FindByID(id string) (*model.Transfer, error)
	CreateTransfer(tx *model.Transfer) error
	UpdateTransferStatus(id string, status string) error
	GetUserWalletID(userID string) (string, error)
	GetWalletByID(walletID string) (*model.Wallet, error)
	UpdateWalletBalance(walletID string, newBalance float64) error
	CreateExternalTransfer(et *model.ExternalTransfer) error
}

type Wallet struct {
	ID      string
	Balance float64
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) FindByID(id string) (*model.Transfer, error) {
	tx := &model.Transfer{}
	query := `
		SELECT id, from_wallet_id, to_wallet_id, amount, status, description, created_at
		FROM transactions
		WHERE id = ? AND type IN ('transfer_in', 'transfer_out')
	`

	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.FromWalletID,
		&tx.ToWalletID,
		&tx.Amount,
		&tx.Status,
		&tx.Description,
		&tx.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return tx, err
}

func (r *transferRepository) CreateTransfer(tx *model.Transfer) error {
	tx.ID = uuid.New().String()
	tx.CreatedAt = time.Now()

	query := `
		INSERT INTO transactions (id, wallet_id, type, amount, balance_before, balance_after, status, reference_id, description, created_at)
		VALUES (?, ?, 'transfer_out', ?, 0, 0, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		tx.ID,
		tx.FromWalletID,
		tx.Amount,
		tx.Status,
		tx.ToWalletID,
		tx.Description,
		tx.CreatedAt,
	)

	return err
}

func (r *transferRepository) UpdateTransferStatus(id string, status string) error {
	query := "UPDATE transactions SET status = ? WHERE id = ?"
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *transferRepository) GetUserWalletID(userID string) (string, error) {
	var walletID string
	query := "SELECT id FROM wallets WHERE user_id = ?"
	err := r.db.QueryRow(query, userID).Scan(&walletID)
	return walletID, err
}

func (r *transferRepository) GetWalletByID(walletID string) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	query := "SELECT id, balance FROM wallets WHERE id = ?"
	err := r.db.QueryRow(query, walletID).Scan(&wallet.ID, &wallet.Balance)
	return wallet, err
}

func (r *transferRepository) UpdateWalletBalance(walletID string, newBalance float64) error {
	query := "UPDATE wallets SET balance = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.Exec(query, newBalance, time.Now(), walletID)
	return err
}

func (r *transferRepository) CreateExternalTransfer(et *model.ExternalTransfer) error {
	et.ID = uuid.New().String()
	et.CreatedAt = time.Now()

	query := `
		INSERT INTO external_transfers (id, transaction_id, provider, recipient_account, recipient_name, amount, fee, status, provider_reference, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		et.ID,
		et.TransactionID,
		et.Provider,
		et.RecipientAccount,
		et.RecipientName,
		et.Amount,
		et.Fee,
		et.Status,
		et.ProviderReference,
		et.CreatedAt,
	)

	return err
}
