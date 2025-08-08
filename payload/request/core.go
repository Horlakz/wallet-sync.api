package request

type WalletFundRequest struct {
	Amount float64 `json:"amount"`
}

type WalletWithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type WalletTransferRequest struct {
	ToAccountNumber string  `json:"to_account_number"`
	Amount          float64 `json:"amount"`
}
