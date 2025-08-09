package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/horlakz/wallet-sync.api/internal/helper"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/shopspring/decimal"
)

type Account struct {
	database.BaseModel

	UserID      *uuid.UUID      `json:"user_id" gorm:"index:idx_user_id;type:uuid"`
	AccountType string          `json:"account_type" gorm:"default:'wallet';not null"`
	Currency    string          `json:"currency" gorm:"type:varchar(10);default:'NGN';not null"`
	Balance     decimal.Decimal `json:"balance" gorm:"not null; type:decimal(32,2)"`
	Number      string          `json:"number" gorm:"type:varchar(20);not null"`
}

type TransactionType string

const (
	Credit               TransactionType = "credit"
	Debit                TransactionType = "debit"
	TransactionPending   string          = "pending"
	TransactionCompleted string          = "completed"
	TransactionFailed    string          = "failed"
)

type Transaction struct {
	database.BaseModel

	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid"`
	Reference   string          `json:"reference" gorm:"not null"`
	Type        TransactionType `json:"type" gorm:"type:enum('credit','debit');not null"`
	Status      string          `json:"status" gorm:"type:enum('pending','completed','failed');default:'pending';not null"`
	Amount      decimal.Decimal `json:"amount" gorm:"not null, type:decimal(32,2)"`
	Currency    string          `json:"currency" gorm:"type:varchar(10);default:'NGN';not null"`
	Description string          `json:"description" gorm:"type:varchar(255);not null"`
	Metadata    datatypes.JSON
}

type LedgerEntry struct {
	database.BaseModel

	UserID        uuid.UUID       `json:"user_id" gorm:"type:uuid"`
	AccountID     uuid.UUID       `json:"account_id" gorm:"type:uuid"`
	TransactionID uuid.UUID       `json:"transaction_id" gorm:"type:uuid"`
	EntryType     string          `json:"entry_type" gorm:"type:enum('debit','credit');not null"`
	Amount        decimal.Decimal `json:"amount" gorm:"not null; type:decimal(32,2)"`
	Description   string          `json:"description" gorm:"type:varchar(255);not null"`
}

type ReconciliationLog struct {
	database.BaseModel

	AccountID       uuid.UUID       `json:"account_id"  gorm:"type:uuid; not null"`
	ComputedBalance decimal.Decimal `json:"computed_balance" gorm:"type:decimal(32,2);not null"`
	StoredBalance   decimal.Decimal `json:"stored_balance" gorm:"type:decimal(32,2);not null"`
	Discrepancy     decimal.Decimal `json:"discrepancy" gorm:"type:decimal(32, 2);generated always as (actual_balance - computed_balance) stored"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a unique reference for the transaction
	t.Reference = helper.GenerateRandomAlphaNumeric(12)

	t.BaseModel.BeforeCreate(tx)

	return nil
}
