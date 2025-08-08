package seed

import (
	"fmt"

	"github.com/horlakz/wallet-sync.api/internal/helper"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
)

type SeederInterface interface {
	Seed()
}

type seeder struct {
	dbConn database.DatabaseInterface
}

func NewSeeder(dbConn database.DatabaseInterface) SeederInterface {
	return &seeder{dbConn: dbConn}
}

func (s *seeder) Seed() {
	s.SeedAdmin()
}

func (s *seeder) SeedAdmin() {
	hashing := helper.NewHashing()
	users := []struct {
		Name  string
		Email string
	}{
		{"Admin User", "admin@wallet-sync.com"},
		{"Test User", "test@wallet-sync.com"},
	}

	hashedPassword, err := hashing.HashPassword("Pa$$w0rd!")
	if err != nil {
		fmt.Println("Failed to hash password:", err)
		return
	}

	for _, userInfo := range users {
		if exists := s.dbConn.Connection().Where("email = ?", userInfo.Email).First(&model.User{}).RowsAffected > 0; exists {
			fmt.Printf("%s already exists in the database. Skipping seeding...\n", userInfo.Name)
			continue
		}

		user := model.User{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Password: hashedPassword,
		}

		if err := s.dbConn.Connection().Create(&user).Error; err != nil {
			fmt.Printf("Failed to create %s: %v\n", userInfo.Name, err)
		} else {
			fmt.Printf("%s created successfully.\n", userInfo.Name)
		}
	}

}
