package response

import (
	"github.com/horlakz/wallet-sync.api/dto"
)

type LoginResponse struct {
	Response

	Data dto.LoginResponseDTO `json:"data"`
}
