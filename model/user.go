package model

import (
	"strings"

	"github.com/horlakz/wallet-sync.api/internal/helper"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type User struct {
	database.BaseModel

	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// convert email to lowercase before saving
	u.Email = strings.ToLower(u.Email)

	// create account model
	account := &Account{
		UserID:      &u.ID,
		AccountType: "wallet",
		Currency:    "NGN",
		Balance:     decimal.NewFromFloat(0.00),
		Number:      helper.GenerateAccountNumber(),
	}

	if err := tx.Create(account).Error; err != nil {
		return err
	}

	return nil
}
