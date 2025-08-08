package core_repository

import (
	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
)

type LedgerEntryRepository interface {
	CreateLedgerEntry(entry *model.LedgerEntry) error
	GetLedgerEntriesByUserID(userID string) ([]model.LedgerEntry, error)
	UpdateLedgerEntry(entry *model.LedgerEntry) error
	GetLedgerEntryByTransactionID(transactionID uuid.UUID) (*model.LedgerEntry, error)
}

type ledgerEntryRepository struct {
	db database.DatabaseInterface
}

func NewLedgerEntryRepository(db database.DatabaseInterface) LedgerEntryRepository {
	return &ledgerEntryRepository{db: db}
}

func (r *ledgerEntryRepository) CreateLedgerEntry(entry *model.LedgerEntry) error {
	return r.db.Connection().Create(entry).Error
}

func (r *ledgerEntryRepository) GetLedgerEntriesByUserID(userID string) ([]model.LedgerEntry, error) {
	var entries []model.LedgerEntry
	err := r.db.Connection().Where("user_id = ?", userID).Find(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *ledgerEntryRepository) UpdateLedgerEntry(entry *model.LedgerEntry) error {
	return r.db.Connection().Save(entry).Error
}

func (r *ledgerEntryRepository) GetLedgerEntryByTransactionID(transactionID uuid.UUID) (*model.LedgerEntry, error) {
	var entry model.LedgerEntry
	err := r.db.Connection().Where("transaction_id = ?", transactionID).First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}
