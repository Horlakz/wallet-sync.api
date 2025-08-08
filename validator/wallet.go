package validator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/horlakz/wallet-sync.api/payload/request"
)

type WalletValidator struct {
	Validator[request.WalletFundRequest]
}

func (validator *WalletValidator) FundValidate(fundReq request.WalletFundRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&fundReq,
		validation.Field(&fundReq.Amount, validation.Required, validation.Min(1)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *WalletValidator) WithdrawValidate(withdrawReq request.WalletWithdrawRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&withdrawReq,
		validation.Field(&withdrawReq.Amount, validation.Required, validation.Min(1)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *WalletValidator) TransferValidate(transferReq request.WalletTransferRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&transferReq,
		validation.Field(&transferReq.ToAccountNumber, validation.Required, validation.Length(10, 10)),
		validation.Field(&transferReq.Amount, validation.Required, validation.Min(1)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
