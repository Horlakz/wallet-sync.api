package dto

import (
	"github.com/shopspring/decimal"
)

type WalletDetailsDto struct {
	Balance       decimal.Decimal `json:"balance"`
	AccountNumber string          `json:"account_number"`
}

type TransactionDto struct {
	Reference   string          `json:"reference"`
	Type        string          `json:"type"`
	Status      string          `json:"status"`
	Amount      decimal.Decimal `json:"amount"`
	Currency    string          `json:"currency"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}
