package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vitaly06/portfolio-rest-api/internal/config"
	"github.com/vitaly06/portfolio-rest-api/internal/domain"
)

type AIService struct {
	cfg *config.Config
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{cfg: cfg}
}

func (s *AIService) AnalyzeAndReply(comment string) domain.AIResult {
	// Если ключ не задан, сразу уходим в fallback
	if s.cfg.OpenAIEkey == "" {
		return s.getFallbackResult()
	}

	prompt := fmt.Sprintf(
		"Ты — AI-ассистент. Твоя задача — проанализировать комментарий пользователя и сгенерировать автоответ. "+
			"Определи тональность комментария (доступно только 3 варианта: 'positive', 'neutral', 'negative'). "+
			"Напиши вежливый, профессиональный краткий ответ от лица Виталия. "+
			"Верни ответ СТРОГО в формате JSON без markdown-разметки (без ```json), содержащий поля 'sentiment' и 'reply'. "+
			"Комментарий: \"%s\"", comment,
	)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	})

	req, err := http.NewRequest("POST", "[https://api.openai.com/v1/chat/completions](https://api.openai.com/v1/chat/completions)", bytes.NewBuffer(requestBody))
	if err != nil {
		return s.getFallbackResult()
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.OpenAIEkey)

	client := &http.Client{Timeout: 7 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return s.getFallbackResult()
	}
	defer resp.Body.Close()

	var openAIResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil || len(openAIResp.Choices) == 0 {
		return s.getFallbackResult()
	}

	var result domain.AIResult
	// Пытаемся распарсить то, что вернула нейросеть в контенте
	err = json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &result)
	if err != nil {
		return s.getFallbackResult()
	}

	return result
}

// Graceful Fallback
func (s *AIService) getFallbackResult() domain.AIResult {
	return domain.AIResult{
		Sentiment: "neutral",
		Reply:     "Здравствуйте! Спасибо за ваше обращение. Я получил ваше сообщение и свяжусь с вами в ближайшее время для обсуждения деталей.",
	}
}
