package service

import (
	"github.com/horlakz/wallet-sync.api/model"
	core_repository "github.com/horlakz/wallet-sync.api/repository/core"
)

type ReconciliationService interface {
	ReconcileTransactions() error
}

type reconciliationService struct {
	ledgerEntryRepo       core_repository.LedgerEntryRepository
	transactionRepo       core_repository.TransactionRepository
	accountRepo           core_repository.AccountRepository
	reconciliationLogRepo core_repository.ReconciliationLogRepository
}

func NewReconciliationService(
	ledgerEntryRepo core_repository.LedgerEntryRepository,
	transactionRepo core_repository.TransactionRepository,
	accountRepo core_repository.AccountRepository,
	reconciliationLogRepo core_repository.ReconciliationLogRepository,
) ReconciliationService {
	return &reconciliationService{
		ledgerEntryRepo:       ledgerEntryRepo,
		transactionRepo:       transactionRepo,
		accountRepo:           accountRepo,
		reconciliationLogRepo: reconciliationLogRepo,
	}
}

func (s *reconciliationService) ReconcileTransactions() error {
	// Fetch all transactions
	transactions, err := s.transactionRepo.GetAllTransactions()
	if err != nil {
		return err
	}

	for _, tx := range transactions {
		// Check if transaction has corresponding ledger entry
		_, err := s.ledgerEntryRepo.GetLedgerEntryByTransactionID(tx.ID)
		if err != nil {
			return err
		}

		// if ledgerEntry == nil {
		// 	// If missing, create ledger entry
		// 	newEntry := model.LedgerEntry{
		// 		TransactionID: tx.ID,
		// 		AccountID:     *tx.AccountID,
		// 		Amount:        tx.Amount,
		// 		Type:          tx.Type,
		// 		CreatedAt:     tx.CreatedAt,
		// 	}
		// 	if err := s.ledgerEntryRepo.CreateLedgerEntry(newEntry); err != nil {
		// 		return err
		// 	}
		// }

		// Optionally, update account balance if needed
		if err != nil {
			return err
		}
		// Example: update balance based on transaction type
		switch tx.Type {
		case "credit":
			if err := s.accountRepo.UpdateAccountBalance(tx.UserID, tx.Amount); err != nil {
				return err
			}
		case "debit":
			if err := s.accountRepo.UpdateAccountBalance(tx.UserID, tx.Amount.Neg()); err != nil {
				return err
			}
		}
	}

	log := model.ReconciliationLog{
		// Populate log fields
	}

	if err := s.reconciliationLogRepo.CreateReconciliationLog(&log); err != nil {
		return err
	}

	return nil
}
