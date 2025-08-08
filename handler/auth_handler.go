package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/wallet-sync.api/dto"
	"github.com/horlakz/wallet-sync.api/payload/request"
	"github.com/horlakz/wallet-sync.api/payload/response"
	"github.com/horlakz/wallet-sync.api/service"
	"github.com/horlakz/wallet-sync.api/validator"
)

type authHandler struct {
	authService service.AuthServiceInterface
	validator   validator.AuthValidator
}

type AuthHandlerInterface interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

func NewAuthHandler(authService service.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{authService: authService}
}

func (handler *authHandler) Login(c *fiber.Ctx) error {
	var resp response.LoginResponse

	loginRequest := new(request.LoginRequest)

	if err := c.BodyParser(loginRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if _, err := handler.validator.LoginValidate(*loginRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	token, err := handler.authService.Login(loginRequest.Email, loginRequest.Password)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data.AccessToken = token

	return c.JSON(resp)
}

func (handler *authHandler) Register(c *fiber.Ctx) error {
	var resp response.Response
	var authDto dto.RegisterDTO

	registerRequest := new(request.RegisterRequest)

	if err := c.BodyParser(registerRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.RegisterValidate(*registerRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	authDto.Name = registerRequest.Name
	authDto.Email = registerRequest.Email
	authDto.Password = registerRequest.Password

	if err := handler.authService.Register(authDto); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = "Registration successful"

	return c.JSON(resp)
}
