package usecase

import (
	"github.com/vitaly06/portfolio-rest-api/internal/domain"
	"github.com/vitaly06/portfolio-rest-api/internal/repository"
	"github.com/vitaly06/portfolio-rest-api/pkg/mailer"
)

type ContactUsecase struct {
	repo      *repository.FileRepository
	aiService *AIService
	mailer    *mailer.Mailer
}

func NewContactUsecase(repo *repository.FileRepository, ai *AIService, m *mailer.Mailer) *ContactUsecase {
	return &ContactUsecase{
		repo:      repo,
		aiService: ai,
		mailer:    m,
	}
}

func (u *ContactUsecase) ProcessContactForm(req domain.ContactRequest) (domain.AIResult, error) {
	// 1. Запуск ИИ аналитики (с fallback защитой)
	aiResult := u.aiService.AnalyzeAndReply(req.Comment)

	// 2. Асинхронно отправляем почту, чтобы не тормозить HTTP-ответ
	go func() {
		_ = u.mailer.SendContactEmails(req.Name, req.Email, req.Phone, req.Comment, aiResult.Reply)
	}()

	// 3. Обновление файловой статистики
	_ = u.repo.UpdateMetrics(aiResult.Sentiment)

	return aiResult, nil
}
