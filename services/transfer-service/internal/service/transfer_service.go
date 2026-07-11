package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ikhdamw/e-wallet/transfer-service/internal/model"
	"github.com/ikhdamw/e-wallet/transfer-service/internal/repository"
	"github.com/ikhdamw/e-wallet/transfer-service/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TransferService interface {
	InternalTransfer(userID string, req *model.InternalTransferRequest) (*model.TransferResponse, error)
	ExternalTransfer(userID string, req *model.ExternalTransferRequest) (*model.TransferResponse, error)
	GetStatus(transferID string) (*model.Transfer, error)
}

type transferService struct {
	transferRepo repository.TransferRepository
	rabbitMQ     *amqp.Connection
	config       *config.Config
}

func NewTransferService(transferRepo repository.TransferRepository, rabbitMQ *amqp.Connection, cfg *config.Config) TransferService {
	return &transferService{
		transferRepo: transferRepo,
		rabbitMQ:     rabbitMQ,
		config:       cfg,
	}
}

func (s *transferService) InternalTransfer(userID string, req *model.InternalTransferRequest) (*model.TransferResponse, error) {
	// Get sender wallet
	senderWalletID, err := s.transferRepo.GetUserWalletID(userID)
	if err != nil {
		return nil, errors.New("sender wallet not found")
	}

	// Get sender wallet details
	senderWallet, err := s.transferRepo.GetWalletByID(senderWalletID)
	if err != nil {
		return nil, errors.New("sender wallet not found")
	}

	// Check balance
	if senderWallet.Balance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	// TODO: Get recipient wallet by email
	// For now, we'll create a pending transfer
	transfer := &model.Transfer{
		FromWalletID: senderWalletID,
		ToWalletID:   "", // Will be resolved
		Amount:       req.Amount,
		Status:       "pending",
		Description:  req.Description,
	}

	// Create transfer record
	if err := s.transferRepo.CreateTransfer(transfer); err != nil {
		return nil, err
	}

	// Publish to RabbitMQ for async processing
	if err := s.publishTransfer(transfer); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to publish transfer: %v\n", err)
	}

	return &model.TransferResponse{
		TransferID: transfer.ID,
		Status:     transfer.Status,
		Amount:     transfer.Amount,
	}, nil
}

func (s *transferService) ExternalTransfer(userID string, req *model.ExternalTransferRequest) (*model.TransferResponse, error) {
	// Get sender wallet
	senderWalletID, err := s.transferRepo.GetUserWalletID(userID)
	if err != nil {
		return nil, errors.New("sender wallet not found")
	}

	// Get sender wallet details
	senderWallet, err := s.transferRepo.GetWalletByID(senderWalletID)
	if err != nil {
		return nil, errors.New("sender wallet not found")
	}

	// Calculate fee (1%)
	fee := req.Amount * 0.01
	totalAmount := req.Amount + fee

	// Check balance
	if senderWallet.Balance < totalAmount {
		return nil, errors.New("insufficient balance (including fee)")
	}

	// Create transfer record
	transfer := &model.Transfer{
		FromWalletID: senderWalletID,
		ToWalletID:   "",
		Amount:       req.Amount,
		Status:       "pending",
		Description:  req.Description,
	}

	if err := s.transferRepo.CreateTransfer(transfer); err != nil {
		return nil, err
	}

	// Create external transfer record
	externalTransfer := &model.ExternalTransfer{
		TransactionID:    transfer.ID,
		Provider:         req.Provider,
		RecipientAccount: req.RecipientAccount,
		Amount:           req.Amount,
		Fee:              fee,
		Status:           "pending",
	}

	if err := s.transferRepo.CreateExternalTransfer(externalTransfer); err != nil {
		return nil, err
	}

	// Publish to RabbitMQ for async processing
	if err := s.publishExternalTransfer(externalTransfer); err != nil {
		fmt.Printf("Failed to publish external transfer: %v\n", err)
	}

	return &model.TransferResponse{
		TransferID: transfer.ID,
		Status:     transfer.Status,
		Amount:     transfer.Amount,
		Fee:        fee,
	}, nil
}

func (s *transferService) GetStatus(transferID string) (*model.Transfer, error) {
	return s.transferRepo.FindByID(transferID)
}

func (s *transferService) publishTransfer(transfer *model.Transfer) error {
	ch, err := s.rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("transfer_internal", true, false, false, false, nil)
	if err != nil {
		return err
	}

	body := fmt.Sprintf(`{"id":"%s","from":"%s","to":"%s","amount":%.2f}`,
		transfer.ID, transfer.FromWalletID, transfer.ToWalletID, transfer.Amount)

	return ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
	})
}

func (s *transferService) publishExternalTransfer(et *model.ExternalTransfer) error {
	ch, err := s.rabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("transfer_external", true, false, false, false, nil)
	if err != nil {
		return err
	}

	body := fmt.Sprintf(`{"id":"%s","provider":"%s","account":"%s","amount":%.2f,"fee":%.2f}`,
		et.ID, et.Provider, et.RecipientAccount, et.Amount, et.Fee)

	return ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
	})
}
