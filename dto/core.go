package dto

import (
	"github.com/shopspring/decimal"
)

type WalletDetailsDto struct {
	Balance       decimal.Decimal `json:"balance"`
	AccountNumber string          `json:"account_number"`
}
