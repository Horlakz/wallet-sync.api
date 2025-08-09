package core_repository

import (
	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(account *model.Account) error
	GetAccountByUserID(userID uuid.UUID) (*model.Account, error)
	GetWalletAccountByUserID(userID uuid.UUID) (*model.Account, error)
	GetAccountByNumber(accountNumber string) (*model.Account, error)
	UpdateAccountBalance(userID uuid.UUID, amount decimal.Decimal) error
	GetAllAccounts() ([]*model.Account, error)
	WithTx(tx *gorm.DB) AccountRepository
}

type accountRepository struct {
	db database.DatabaseInterface
}

func NewAccountRepository(db database.DatabaseInterface) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) WithTx(tx *gorm.DB) AccountRepository {
	return &accountRepository{db: database.Wrap(tx)}
}

func (r *accountRepository) CreateAccount(account *model.Account) error {
	return r.db.Connection().Create(account).Error
}

func (r *accountRepository) GetAccountByUserID(userID uuid.UUID) (*model.Account, error) {
	var account model.Account
	err := r.db.Connection().Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetWalletAccountByUserID(userID uuid.UUID) (*model.Account, error) {
	var account model.Account
	err := r.db.Connection().Where("user_id = ? AND account_type = ?", userID, "wallet").First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetAccountByNumber(accountNumber string) (*model.Account, error) {
	var account model.Account
	err := r.db.Connection().Where("number = ?", accountNumber).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) UpdateAccountBalance(userID uuid.UUID, amount decimal.Decimal) error {
	var account model.Account
	err := r.db.Connection().Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		return err
	}

	account.Balance = account.Balance.Add(amount)
	return r.db.Connection().Save(&account).Error
}

func (r *accountRepository) GetAllAccounts() ([]*model.Account, error) {
	var accounts []*model.Account
	err := r.db.Connection().Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
