package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/vitaly06/portfolio-rest-api/internal/domain"
	"github.com/vitaly06/portfolio-rest-api/internal/repository"
	"github.com/vitaly06/portfolio-rest-api/internal/usecase"
)

type Handler struct {
	usecase  *usecase.ContactUsecase
	repo     *repository.FileRepository
	validate *validator.Validate
}

func NewHandler(u *usecase.ContactUsecase, r *repository.FileRepository) *Handler {
	return &Handler{
		usecase:  u,
		repo:     r,
		validate: validator.New(),
	}
}

func (h *Handler) HandleContact(c fiber.Ctx) error {
	var req domain.ContactRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный формат JSON"})
	}

	// Валидация данных
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.usecase.ProcessContactForm(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка сервера при обработке"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":   true,
		"sentiment": result.Sentiment,
		"ai_reply":  result.Reply,
	})
}

func (h *Handler) HandleHealth(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "OK"})
}

func (h *Handler) HandleMetrics(c fiber.Ctx) error {
	metrics, err := h.repo.GetMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось прочитать метрики"})
	}
	return c.Status(fiber.StatusOK).JSON(metrics)
}
