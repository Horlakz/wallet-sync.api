package service

import (
	"fmt"

	"github.com/horlakz/wallet-sync.api/internal/config"
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
	logger                *config.Logger
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
		logger:                config.NewLogger(),
	}
}

func (s *reconciliationService) ReconcileTransactions() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered in ReconcileTransactions: %v", r)
			s.logger.Log().Errorf("panic recovered in ReconcileTransactions: %v", r)
			return
		}
	}()

	accounts, err := s.accountRepo.GetAllAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Log().Errorf("panic recovered for account %v: %v", account.ID, r)
				}
			}()

			credits, err := s.ledgerEntryRepo.GetTotalCreditsByAccountID(account.ID)
			if err != nil {
				s.logger.Log().Errorf("error getting credits for account %v: %v", account.ID, err)
				return
			}

			debits, err := s.ledgerEntryRepo.GetTotalDebitsByAccountID(account.ID)
			if err != nil {
				s.logger.Log().Errorf("error getting debits for account %v: %v", account.ID, err)
				return
			}

			computedBalance := credits.Sub(debits)

			if computedBalance.Cmp(account.Balance) != 0 {
				reconciliationLog := &model.ReconciliationLog{
					AccountID:       account.ID,
					ComputedBalance: computedBalance,
					StoredBalance:   account.Balance,
					Discrepancy:     computedBalance.Sub(account.Balance),
				}

				if err := s.reconciliationLogRepo.CreateReconciliationLog(reconciliationLog); err != nil {
					s.logger.Log().Errorf("error creating reconciliation log for account %v: %v", account.ID, err)
					return
				}
			}
		}()
	}

	return nil
}
