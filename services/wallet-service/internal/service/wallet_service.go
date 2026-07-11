package service

import (
	"errors"
	"fmt"

	"github.com/ikhdamw/e-wallet/wallet-service/internal/model"
	"github.com/ikhdamw/e-wallet/wallet-service/internal/repository"
	"github.com/ikhdamw/e-wallet/wallet-service/pkg/config"

	"github.com/redis/go-redis/v9"
)

type WalletService interface {
	GetBalance(userID string) (*model.BalanceResponse, error)
	TopUp(userID string, req *model.TopUpRequest) (*model.Transaction, error)
	GetHistory(userID string, page, limit int) (*model.HistoryResponse, error)
}

type walletService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.WalletRepository
	redis           *redis.Client
	config          *config.Config
}

func NewWalletService(walletRepo, transactionRepo repository.WalletRepository, redis *redis.Client, cfg *config.Config) WalletService {
	return &walletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		redis:           redis,
		config:          cfg,
	}
}

func (s *walletService) GetBalance(userID string) (*model.BalanceResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	return &model.BalanceResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}

func (s *walletService) TopUp(userID string, req *model.TopUpRequest) (*model.Transaction, error) {
	// Find wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Check wallet status
	if wallet.Status != "active" {
		return nil, errors.New("wallet is not active")
	}

	// Calculate new balance
	balanceBefore := wallet.Balance
	balanceAfter := balanceBefore + req.Amount

	// Update balance
	if err := s.walletRepo.UpdateBalance(wallet.ID, balanceAfter); err != nil {
		return nil, err
	}

	// Create transaction
	tx := &model.Transaction{
		WalletID:      wallet.ID,
		Type:          "topup",
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        "completed",
		Description:   fmt.Sprintf("Top up IDR %.2f", req.Amount),
	}

	if err := s.transactionRepo.CreateTransaction(tx); err != nil {
		// Rollback balance
		_ = s.walletRepo.UpdateBalance(wallet.ID, balanceBefore)
		return nil, err
	}

	return tx, nil
}

func (s *walletService) GetHistory(userID string, page, limit int) (*model.HistoryResponse, error) {
	// Find wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Get transactions
	transactions, total, err := s.transactionRepo.GetTransactions(wallet.ID, page, limit)
	if err != nil {
		return nil, err
	}

	return &model.HistoryResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}
