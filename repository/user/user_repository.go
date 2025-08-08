package user_repository

import (
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db database.DatabaseInterface
}

func NewUserRepository(db database.DatabaseInterface) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Connection().Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Connection().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
