package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model

	ID uuid.UUID `gorm:"primaryKey;" json:"id"`
}

func (model *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	uid, err := uuid.NewV7()

	if err != nil {
		return err
	}

	model.ID = uid

	return nil
}
