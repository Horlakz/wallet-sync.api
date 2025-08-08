package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/horlakz/wallet-sync.api/payload/request"
	"github.com/horlakz/wallet-sync.api/payload/response"
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
	userId := c.Locals("userId").(uuid.UUID)
	var resp response.Response

	wallet, err := handler.walletService.GetWalletDetails(userId)
	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Failed to retrieve wallet details"
		return c.Status(resp.Status).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Wallet details retrieved successfully"
	resp.Data = wallet

	return c.Status(resp.Status).JSON(resp)
}

func (handler *walletHandler) Transfer(c *fiber.Ctx) error {
	var transferRequest request.WalletTransferRequest
	var resp response.Response

	if err := c.BodyParser(&transferRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(resp.Status).JSON(resp)
	}

	if _, err := handler.validator.TransferValidate(transferRequest); err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = err.Error()
		return c.Status(resp.Status).JSON(resp)
	}

	userId := c.Locals("userId").(uuid.UUID)

	amountDecimal := decimal.NewFromFloat(transferRequest.Amount)
	handler.walletService.TransferFunds(userId, transferRequest.ToAccountNumber, amountDecimal)

	resp.Status = http.StatusOK
	resp.Message = "Transfer successful"
	return c.Status(resp.Status).JSON(resp)
}

func (handler *walletHandler) Fund(c *fiber.Ctx) error {
	var fundRequest request.WalletFundRequest
	var resp response.Response

	if err := c.BodyParser(&fundRequest); err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "Invalid request"
		return c.Status(resp.Status).JSON(resp)
	}

	if _, err := handler.validator.FundValidate(fundRequest); err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = err.Error()
		return c.Status(resp.Status).JSON(resp)
	}

	userId := c.Locals("userId").(uuid.UUID)

	amountDecimal := decimal.NewFromFloat(fundRequest.Amount)

	if err := handler.walletService.FundWallet(userId, amountDecimal); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(resp.Status).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Wallet funded successfully"
	return c.Status(resp.Status).JSON(resp)
}

func (handler *walletHandler) Withdraw(c *fiber.Ctx) error {
	var withdrawRequest request.WalletWithdrawRequest
	var resp response.Response

	if err := c.BodyParser(&withdrawRequest); err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = "Invalid request"
		return c.Status(resp.Status).JSON(resp)
	}

	if _, err := handler.validator.WithdrawValidate(withdrawRequest); err != nil {
		resp.Status = http.StatusUnprocessableEntity
		resp.Message = err.Error()
		return c.Status(resp.Status).JSON(resp)
	}

	userId := c.Locals("userId").(uuid.UUID)

	amountDecimal := decimal.NewFromFloat(withdrawRequest.Amount)

	if err := handler.walletService.WithdrawFromWallet(userId, amountDecimal); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(resp.Status).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Withdrawal successful"
	return c.Status(resp.Status).JSON(resp)
}
