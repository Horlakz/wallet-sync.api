package core_repository

import (
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
)

type ReconciliationLogRepository interface {
	CreateReconciliationLog(log *model.ReconciliationLog) error
	GetReconciliationLogsByUserID(userID string) ([]model.ReconciliationLog, error)
	UpdateReconciliationLog(log *model.ReconciliationLog) error
}

type reconciliationLogRepository struct {
	db database.DatabaseInterface
}

func NewReconciliationLogRepository(db database.DatabaseInterface) ReconciliationLogRepository {
	return &reconciliationLogRepository{db: db}
}

func (r *reconciliationLogRepository) CreateReconciliationLog(log *model.ReconciliationLog) error {
	return r.db.Connection().Create(log).Error
}

func (r *reconciliationLogRepository) GetReconciliationLogsByUserID(userID string) ([]model.ReconciliationLog, error) {
	var logs []model.ReconciliationLog
	err := r.db.Connection().Where("user_id = ?", userID).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *reconciliationLogRepository) UpdateReconciliationLog(log *model.ReconciliationLog) error {
	return r.db.Connection().Save(log).Error
}
