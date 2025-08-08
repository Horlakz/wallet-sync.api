package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/payload/request"
	"github.com/horlakz/wallet-sync.api/service"
	"github.com/horlakz/wallet-sync.api/validator"
	"github.com/shopspring/decimal"
)

type walletHandler struct {
	walletService service.WalletServiceInterface
	validator     validator.WalletValidator
}

type WalletHandlerInterface interface {
	GetDetails(c *fiber.Ctx) error
	Transfer(c *fiber.Ctx) error
	Fund(c *fiber.Ctx) error
	Withdraw(c *fiber.Ctx) error
}

func NewWalletHandler(walletService service.WalletServiceInterface) WalletHandlerInterface {
	return &walletHandler{walletService: walletService}
}

func (handler *walletHandler) GetDetails(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	wallet, err := handler.walletService.GetWalletDetails(uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(wallet)
}

func (handler *walletHandler) Transfer(c *fiber.Ctx) error {
	var transferRequest request.WalletTransferRequest

	if err := c.BodyParser(&transferRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if _, err := handler.validator.TransferValidate(transferRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	userId := c.Locals("userId").(string)

	uuid, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	amountDecimal := decimal.NewFromFloat(transferRequest.Amount)
	handler.walletService.TransferFunds(uuid, transferRequest.ToAccountNumber, amountDecimal)

	return c.JSON(fiber.Map{"message": "Transfer successful"})
}

func (handler *walletHandler) Fund(c *fiber.Ctx) error {
	var fundRequest request.WalletFundRequest

	if err := c.BodyParser(&fundRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if _, err := handler.validator.FundValidate(fundRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	userId := c.Locals("userId").(string)
	uuid, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	amountDecimal := decimal.NewFromFloat(fundRequest.Amount)

	if err := handler.walletService.FundWallet(uuid, amountDecimal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Wallet funded successfully"})
}

func (handler *walletHandler) Withdraw(c *fiber.Ctx) error {
	var withdrawRequest request.WalletWithdrawRequest

	if err := c.BodyParser(&withdrawRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if _, err := handler.validator.WithdrawValidate(withdrawRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	userId := c.Locals("userId").(string)
	uuid, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	amountDecimal := decimal.NewFromFloat(withdrawRequest.Amount)

	if err := handler.walletService.WithdrawFromWallet(uuid, amountDecimal); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Withdrawal successful"})
}
