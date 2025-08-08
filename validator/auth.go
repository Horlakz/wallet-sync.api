package validator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/horlakz/wallet-sync.api/payload/request"
)

type AuthValidator struct {
	Validator[request.RegisterRequest]
}

func (validator *AuthValidator) LoginValidate(loginReq request.LoginRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&loginReq,
		validation.Field(&loginReq.Email, validation.Required, validation.Length(3, 32), is.Email),
		validation.Field(&loginReq.Password, validation.Required, validation.Length(3, 32)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *AuthValidator) RegisterValidate(registerDto request.RegisterRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&registerDto,
		validation.Field(&registerDto.Name, validation.Required, validation.Length(3, 32)),
		validation.Field(&registerDto.Email, validation.Required, validation.Length(3, 32)),
		validation.Field(&registerDto.Password, validation.Required, validation.Length(3, 32)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
