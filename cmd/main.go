package main

import (
	"log"
	"os"

	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/vitaly06/portfolio-rest-api/internal/config"
	"github.com/vitaly06/portfolio-rest-api/internal/deliviry/http"
	"github.com/vitaly06/portfolio-rest-api/internal/repository"
	"github.com/vitaly06/portfolio-rest-api/internal/usecase"
	"github.com/vitaly06/portfolio-rest-api/pkg/mailer"
)

func main() {
	cfg := config.LoadConfig()

	// Создаем папки для хранения данных, если их нет
	_ = os.MkdirAll("data", 0755)

	// Инициализация слоев
	repo := repository.NewFileRepository("data/app.log", "data/stats.json")
	aiService := usecase.NewAIService(cfg)
	mailService := mailer.NewMailer(cfg)
	contactUsecase := usecase.NewContactUsecase(repo, aiService, mailService)
	handler := http.NewHandler(contactUsecase, repo)

	app := fiber.New(fiber.Config{
		AppName: "Developer Portfolio API v1.0",
	})

	// Глобальные Middleware
	app.Use(cors.New())
	app.Use(http.NewLoggerMiddleware(repo))

	// Ограничение частоты запросов
	app.Use("/api/contact", limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Слишком много запросов. Пожалуйста, подождите минуту.",
			})
		},
	}))

	app.Post("/api/contact", handler.HandleContact)
	app.Get("/api/health", handler.HandleHealth)
	app.Get("/api/metrics", handler.HandleMetrics)

	log.Printf("Сервер успешно запущен на порту %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка старта сервера: %v", err)
	}
}
