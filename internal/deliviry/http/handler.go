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

// HandleContact godoc
// @Summary      Отправка формы обратной связи
// @Description  Принимает данные формы, анализирует тональность через AI, шлет email и возвращает автоответ
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        request body domain.ContactRequest true "Данные контакта"
// @Success      200 {object} domain.ContactResponse "Успешный ответ с AI-автоответом"
// @Failure      400 {object} domain.ErrorResponse "Ошибка валидации или неверный JSON"
// @Failure      429 {object} domain.ErrorResponse "Too many requests (спам-защита)"
// @Failure      500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /api/contact [post]
func (h *Handler) HandleContact(c fiber.Ctx) error {
	var req domain.ContactRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Некорректный формат JSON"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Error: err.Error()})
	}

	result, err := h.usecase.ProcessContactForm(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Ошибка сервера при обработке"})
	}

	return c.Status(fiber.StatusOK).JSON(domain.ContactResponse{
		Success:   true,
		Sentiment: result.Sentiment,
		AiReply:   result.Reply,
	})
}

// HandleHealth godoc
// @Summary      Проверка статуса сервиса
// @Description  Возвращает текущий статус сервера (health check)
// @Tags         system
// @Produce      json
// @Success      200 {object} domain.HealthResponse "Статус OK"
// @Router       /api/health [get]
func (h *Handler) HandleHealth(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(domain.HealthResponse{Status: "OK"})
}

// HandleMetrics godoc
// @Summary      Статистика обращений
// @Description  Возвращает общие метрики и тональность обращений из файла
// @Tags         system
// @Produce      json
// @Success      200 {object} domain.Metrics "Объект статистики"
// @Failure      500 {object} domain.ErrorResponse "Не удалось прочитать метрики"
// @Router       /api/metrics [get]
func (h *Handler) HandleMetrics(c fiber.Ctx) error {
	metrics, err := h.repo.GetMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Не удалось прочитать метрики"})
	}
	return c.Status(fiber.StatusOK).JSON(metrics)
}
