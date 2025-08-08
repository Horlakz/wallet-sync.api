package core_repository

import (
	"github.com/horlakz/wallet-sync.api/dto"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/model"
)

type Pageable struct {
	Page          int    `json:"page"`
	Size          int    `json:"size"`
	SortBy        string `json:"sort_by"`
	SortDirection string `json:"sort_dir"`
	Type          string `json:"type"`
	Status        string `json:"status"`
}

type Pagination struct {
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}

type TransactionRepository interface {
	CreateTransaction(transaction *model.Transaction) error
	GetTransactionByReference(reference string) (*model.Transaction, error)
	UpdateTransactionStatus(reference string, status string) error
	FindTransactionsByUserID(userID string, pageable Pageable) ([]dto.TransactionDto, Pagination, error)
	GetAllTransactions() ([]model.Transaction, error)
}

type transactionRepository struct {
	db database.DatabaseInterface
}

func NewTransactionRepository(db database.DatabaseInterface) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction *model.Transaction) error {
	return r.db.Connection().Create(transaction).Error
}

func (r *transactionRepository) GetTransactionByReference(reference string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Connection().Where("reference = ?", reference).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateTransactionStatus(reference string, status string) error {
	var transaction model.Transaction
	err := r.db.Connection().Where("reference = ?", reference).First(&transaction).Error
	if err != nil {
		return err
	}

	transaction.Status = status
	return r.db.Connection().Save(&transaction).Error
}

func (r *transactionRepository) GetTransactionsByUserID(userID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Connection().Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindTransactionsByUserID(userID string, pageable Pageable) ([]dto.TransactionDto, Pagination, error) {
	var transactions []dto.TransactionDto
	var totalItems int64

	query := r.db.Connection().Model(&model.Transaction{}).Where("user_id = ?", userID)

	if pageable.Type != "" {
		query = query.Where("type = ?", pageable.Type)
	}

	if pageable.Status != "" {
		query = query.Where("status = ?", pageable.Status)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, Pagination{}, err
	}

	if pageable.SortBy != "" {
		query = query.Order(pageable.SortBy + " " + pageable.SortDirection)
	}

	offset := (pageable.Page - 1) * pageable.Size
	query = query.Offset(offset).Limit(pageable.Size)

	if err := query.Find(&transactions).Error; err != nil {
		return nil, Pagination{}, err
	}

	totalPages := (totalItems + int64(pageable.Size) - 1) / int64(pageable.Size)

	return transactions, Pagination{
		CurrentPage: int64(pageable.Page),
		TotalPages:  totalPages,
		TotalItems:  totalItems,
	}, nil
}

func (r *transactionRepository) GetAllTransactions() ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Connection().Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
