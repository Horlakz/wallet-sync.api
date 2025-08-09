package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/dto"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
	core_repository "github.com/horlakz/wallet-sync.api/repository/core"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletServiceInterface interface {
	FundWallet(userID uuid.UUID, amount decimal.Decimal) error
	WithdrawFromWallet(userID uuid.UUID, amount decimal.Decimal) error
	GetWalletDetails(userID uuid.UUID) (dto.WalletDetailsDto, error)
	TransferFunds(fromUserID uuid.UUID, toAccountNumber string, amount decimal.Decimal) error
}

type walletService struct {
	accountRepo     core_repository.AccountRepository
	transactionRepo core_repository.TransactionRepository
	ledgerEntryRepo core_repository.LedgerEntryRepository
	db              database.DatabaseInterface
}

func NewWalletService(
	accountRepo core_repository.AccountRepository,
	transactionRepo core_repository.TransactionRepository,
	ledgerEntryRepo core_repository.LedgerEntryRepository,
	db database.DatabaseInterface,
) WalletServiceInterface {
	return &walletService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		ledgerEntryRepo: ledgerEntryRepo,
		db:              db,
	}
}

func (s *walletService) FundWallet(userID uuid.UUID, amount decimal.Decimal) error { // use a transaction and rollback if any step fails
	return s.TxHelper(func(
		accountRepo core_repository.AccountRepository,
		transactionRepo core_repository.TransactionRepository,
		ledgerEntryRepo core_repository.LedgerEntryRepository,
	) error {

		account, err := accountRepo.GetAccountByUserID(userID)
		if err != nil {
			return err
		}

		// Update account balance
		if err := accountRepo.UpdateAccountBalance(userID, amount); err != nil {
			return err
		}

		// Create a transaction record
		transaction := model.Transaction{
			UserID:      *account.UserID,
			Type:        model.Credit,
			Status:      model.TransactionCompleted,
			Amount:      amount,
			Currency:    "NGN",
			Description: "Wallet funding",
		}

		if err := transactionRepo.CreateTransaction(&transaction); err != nil {
			return err
		}

		// Create a ledger entry
		ledgerEntry := &model.LedgerEntry{
			UserID:        userID,
			AccountID:     account.ID,
			TransactionID: transaction.ID,
			EntryType:     "credit",
			Amount:        amount,
			Description:   "Wallet funding",
		}

		if err := ledgerEntryRepo.CreateLedgerEntry(ledgerEntry); err != nil {
			return err
		}

		return nil
	})
}

func (s *walletService) WithdrawFromWallet(userID uuid.UUID, amount decimal.Decimal) error {
	return s.TxHelper(func(
		accountRepo core_repository.AccountRepository,
		transactionRepo core_repository.TransactionRepository,
		ledgerEntryRepo core_repository.LedgerEntryRepository,
	) error {

		account, err := accountRepo.GetAccountByUserID(userID)
		if err != nil {
			return err
		}

		// Check if the account has sufficient balance
		if account.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// Update account balance
		if err := accountRepo.UpdateAccountBalance(userID, amount.Neg()); err != nil {
			return err
		}

		// Create a transaction record
		transaction := &model.Transaction{
			UserID:      *account.UserID,
			Type:        model.Debit,
			Status:      model.TransactionCompleted,
			Amount:      amount,
			Currency:    "NGN",
			Description: "Wallet withdrawal",
		}

		if err := transactionRepo.CreateTransaction(transaction); err != nil {
			return err
		}

		// Create a ledger entry
		ledgerEntry := &model.LedgerEntry{
			UserID:        userID,
			AccountID:     account.ID,
			TransactionID: transaction.ID,
			EntryType:     "debit",
			Amount:        amount,
			Description:   "Wallet withdrawal",
		}

		if err := ledgerEntryRepo.CreateLedgerEntry(ledgerEntry); err != nil {
			return err
		}

		return nil
	})

}

func (s *walletService) GetWalletDetails(userID uuid.UUID) (dto.WalletDetailsDto, error) {
	account, err := s.accountRepo.GetWalletAccountByUserID(userID)
	if err != nil {
		return dto.WalletDetailsDto{}, err
	}

	return dto.WalletDetailsDto{
		Balance:       account.Balance,
		AccountNumber: account.Number,
	}, nil
}

func (s *walletService) TransferFunds(fromUserID uuid.UUID, toAccountNumber string, amount decimal.Decimal) error {

	return s.TxHelper(func(
		accountRepo core_repository.AccountRepository,
		transactionRepo core_repository.TransactionRepository,
		ledgerEntryRepo core_repository.LedgerEntryRepository,
	) error {

		fromAccount, err := accountRepo.GetWalletAccountByUserID(fromUserID)
		if err != nil {
			return err
		}

		toAccount, err := accountRepo.GetAccountByNumber(toAccountNumber)
		if err != nil {
			return err
		}

		// check if wallet belongs to self
		if *toAccount.UserID == fromUserID {
			return errors.New("cannot transfer to your own account")
		}

		// Check if the from account has sufficient balance
		if fromAccount.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// Update the balances of both accounts
		if err := accountRepo.UpdateAccountBalance(*fromAccount.UserID, amount.Neg()); err != nil {
			return err
		}
		if err := accountRepo.UpdateAccountBalance(*toAccount.UserID, amount); err != nil {
			return err
		}

		// Create a transaction record for the sender
		senderTransaction := &model.Transaction{
			UserID:      *fromAccount.UserID,
			Type:        model.Debit,
			Status:      model.TransactionCompleted,
			Amount:      amount,
			Currency:    "NGN",
			Description: "Transfer to " + toAccount.Number,
		}

		if err := transactionRepo.CreateTransaction(senderTransaction); err != nil {
			return err
		}

		// Create a ledger entry for the sender
		senderLedgerEntry := &model.LedgerEntry{
			UserID:        *fromAccount.UserID,
			AccountID:     fromAccount.ID,
			TransactionID: senderTransaction.ID,
			EntryType:     "debit",
			Amount:        amount,
			Description:   "Transfer to " + toAccount.Number,
		}

		if err := ledgerEntryRepo.CreateLedgerEntry(senderLedgerEntry); err != nil {
			return err
		}

		// Create a transaction record for the receiver
		receiverTransaction := &model.Transaction{
			UserID:      *toAccount.UserID,
			Type:        model.Credit,
			Status:      model.TransactionCompleted,
			Amount:      amount,
			Currency:    "NGN",
			Description: "Transfer from " + fromAccount.Number,
		}

		if err := transactionRepo.CreateTransaction(receiverTransaction); err != nil {
			return err
		}

		// Create a ledger entry for the receiver
		receiverLedgerEntry := &model.LedgerEntry{
			UserID:        *toAccount.UserID,
			AccountID:     toAccount.ID,
			TransactionID: receiverTransaction.ID,
			EntryType:     "credit",
			Amount:        amount,
			Description:   "Transfer from " + fromAccount.Number,
		}

		if err := ledgerEntryRepo.CreateLedgerEntry(receiverLedgerEntry); err != nil {
			return err
		}

		return nil
	})
}

// TxHelper wraps a function in a DB transaction and injects repository instances with the transaction context.
func (s *walletService) TxHelper(fn func(
	accountRepo core_repository.AccountRepository,
	transactionRepo core_repository.TransactionRepository,
	ledgerEntryRepo core_repository.LedgerEntryRepository,
) error) error {
	return s.db.Connection().Transaction(func(tx *gorm.DB) error {
		accountRepoTx := s.accountRepo.WithTx(tx)
		transactionRepoTx := s.transactionRepo.WithTx(tx)
		ledgerEntryRepoTx := s.ledgerEntryRepo.WithTx(tx)
		return fn(accountRepoTx, transactionRepoTx, ledgerEntryRepoTx)
	})
}
