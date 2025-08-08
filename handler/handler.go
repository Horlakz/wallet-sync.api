package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/payload/response"
	core_repository "github.com/horlakz/wallet-sync.api/repository/core"
)

func GetUserId(c *fiber.Ctx) uuid.UUID {
	userId := c.Locals("userId").(uuid.UUID)

	return userId
}

func Index(c *fiber.Ctx) error {

	var resp response.Response

	var about struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Author  string `json:"author"`
	}

	about.Name = "Wallet Sync API"
	about.Version = "0.0.1"
	about.Author = "Horlakz"

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"about": about}

	return c.JSON(resp)
}

func NotFound(c *fiber.Ctx) error {
	var resp response.Response

	resp.Status = http.StatusNotFound
	resp.Message = "Route not found"

	return c.Status(http.StatusNotFound).JSON(resp)
}

func GetLogs(c *fiber.Ctx, db database.DatabaseInterface) error {
	var resp response.Response

	key := c.Params("key")
	size := c.Query("size", "50")

	intSize, err := strconv.Atoi(size)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()

		return c.JSON(resp)
	}

	logs, err := db.Cache().Get(key, int64(intSize))
	if err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = err.Error()

		return c.JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"logs": logs}

	return c.JSON(resp)
}

func GeneratePageable(context *fiber.Ctx) (pageable core_repository.Pageable) {

	pageable.Page = 1
	pageable.Size = 20
	pageable.SortBy = "created_at"
	pageable.SortDirection = "desc"
	pageable.Type = ""
	pageable.Status = ""

	size, err := strconv.Atoi(context.Query("size", "0"))
	if (size > 0) && err == nil {
		pageable.Size = size
	}

	page, err := strconv.Atoi(context.Query("page", "1"))
	if (page > 0) && err == nil {
		pageable.Page = page
	}

	orderBy := context.Query("sort_by", "")
	if orderBy != "" {
		pageable.SortBy = orderBy
	}

	sortDir := context.Query("sort_dir", "")
	if sortDir != "" {
		pageable.SortBy = sortDir
	}

	_type := context.Query("type", "")
	if _type != "" {
		pageable.Type = _type
	}

	status := context.Query("status", "")
	if status != "" {
		pageable.Status = status
	}

	return pageable
}
