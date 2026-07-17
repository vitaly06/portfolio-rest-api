package http

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/vitaly06/portfolio-rest-api/internal/domain"
	"github.com/vitaly06/portfolio-rest-api/internal/repository"
)

func NewLoggerMiddleware(repo *repository.FileRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next() // Выполняем запрос

		status := c.Response().StatusCode()
		entry := domain.LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Method:    c.Method(),
			Path:      c.Path(),
			IP:        c.IP(),
			Status:    status,
		}

		// Записываем лог в файл асинхронно
		go func() {
			_ = repo.WriteLog(entry)
		}()

		return err
	}
}
