package service

import (
	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/model"
	core_repository "github.com/horlakz/wallet-sync.api/repository/core"
)

type TransactionServiceInterface interface {
	FindTransactionsByUserID(userID uuid.UUID, pageable core_repository.Pageable) ([]model.Transaction, core_repository.Pagination, error)
}

type transactionService struct {
	transactionRepo core_repository.TransactionRepository
}

func NewTransactionService(transactionRepo core_repository.TransactionRepository) TransactionServiceInterface {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) FindTransactionsByUserID(userID uuid.UUID, pageable core_repository.Pageable) ([]model.Transaction, core_repository.Pagination, error) {
	transactions, pagination, err := s.transactionRepo.FindTransactionsByUserID(userID.String(), pageable)
	if err != nil {
		return nil, core_repository.Pagination{}, err
	}

	return transactions, pagination, nil
}
