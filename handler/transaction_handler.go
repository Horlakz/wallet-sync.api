package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/service"
)

type transactionHandler struct {
	transactionService service.TransactionServiceInterface
}

type TransactionHandlerInterface interface {
	GetTransactionHistory(c *fiber.Ctx) error
}

func NewTransactionHandler(transactionService service.TransactionServiceInterface) TransactionHandlerInterface {
	return &transactionHandler{transactionService: transactionService}
}

func (handler *transactionHandler) GetTransactionHistory(c *fiber.Ctx) error {
	userId := c.Locals("userId").(uuid.UUID)
	pageable := GeneratePageable(c)

	transactions, pagination, err := handler.transactionService.FindTransactionsByUserID(userId, pageable)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"pagination":   pagination,
	})
}
