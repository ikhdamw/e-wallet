package repository

import (
	"database/sql"
	"time"

	"github.com/ikhdamw/e-wallet/wallet-service/internal/model"
	"github.com/google/uuid"
)

type WalletRepository interface {
	FindByUserID(userID string) (*model.Wallet, error)
	UpdateBalance(walletID string, newBalance float64) error
	CreateTransaction(tx *model.Transaction) error
	GetTransactions(walletID string, page, limit int) ([]model.Transaction, int, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) FindByUserID(userID string) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	query := `
		SELECT id, user_id, balance, currency, status, created_at, updated_at
		FROM wallets
		WHERE user_id = ?
	`

	err := r.db.QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.Status,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return wallet, err
}

func (r *walletRepository) UpdateBalance(walletID string, newBalance float64) error {
	query := `
		UPDATE wallets
		SET balance = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, newBalance, time.Now(), walletID)
	return err
}

func (r *walletRepository) CreateTransaction(tx *model.Transaction) error {
	tx.ID = uuid.New().String()
	tx.CreatedAt = time.Now()

	query := `
		INSERT INTO transactions (id, wallet_id, type, amount, balance_before, balance_after, status, reference_id, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		tx.ID,
		tx.WalletID,
		tx.Type,
		tx.Amount,
		tx.BalanceBefore,
		tx.BalanceAfter,
		tx.Status,
		tx.ReferenceID,
		tx.Description,
		tx.CreatedAt,
	)

	return err
}

func (r *walletRepository) GetTransactions(walletID string, page, limit int) ([]model.Transaction, int, error) {
	var transactions []model.Transaction
	var total int

	// Count total
	countQuery := "SELECT COUNT(*) FROM transactions WHERE wallet_id = ?"
	err := r.db.QueryRow(countQuery, walletID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get transactions
	offset := (page - 1) * limit
	query := `
		SELECT id, wallet_id, type, amount, balance_before, balance_after, status, reference_id, description, created_at
		FROM transactions
		WHERE wallet_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, walletID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var tx model.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.WalletID,
			&tx.Type,
			&tx.Amount,
			&tx.BalanceBefore,
			&tx.BalanceAfter,
			&tx.Status,
			&tx.ReferenceID,
			&tx.Description,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, total, nil
}
