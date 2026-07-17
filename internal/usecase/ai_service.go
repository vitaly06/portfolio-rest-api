package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	// Берем  Yandex API Key
	apiKey := s.cfg.YandexKey
	if apiKey == "" {
		fmt.Println("[YANDEX ERROR] API ключ пуст")
		return s.getFallbackResult()
	}

	folderID := "b1gk0u23tfejbe5pb3jh"

	systemPrompt := "Ты — AI-ассистент бэкенд-разработчика Виталия. Твоя задача — проанализировать комментарий пользователя и сгенерировать автоответ. " +
		"Определи тональность комментария (доступно только 3 варианта: 'positive', 'neutral', 'negative'). " +
		"Напиши вежливый, профессиональный краткий ответ от лица Виталия. " +
		"Верни ответ СТРОГО в формате JSON без markdown-разметки, содержащий поля 'sentiment' и 'reply'. " +
		"Пример структуры: {\"sentiment\": \"positive\", \"reply\": \"текст\"}."

	userPrompt := fmt.Sprintf("Комментарий для анализа: \"%s\"", comment)

	// Формируем тело запроса
	requestBody, _ := json.Marshal(map[string]interface{}{
		"modelUri": fmt.Sprintf("gpt://%s/yandexgpt-lite/latest", folderID),
		"completionOptions": map[string]interface{}{
			"stream":      false,
			"temperature": 0.3,
			"maxTokens":   2000,
		},
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"text": systemPrompt,
			},
			{
				"role": "user",
				"text": userPrompt,
			},
		},
	})

	apiURL := "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("[YANDEX ERROR] Ошибка создания запроса: %v\n", err)
		return s.getFallbackResult()
	}

	// Авторизация в Yandex Cloud через API-ключ
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", apiKey))

	client := &http.Client{Timeout: 7 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[YANDEX ERROR] Сетевая ошибка: %v\n", err)
		return s.getFallbackResult()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Printf("[YANDEX ERROR] Яндекс вернул статус %d. Ответ: %s\n", resp.StatusCode, buf.String())
		return s.getFallbackResult()
	}

	// Структура ответа
	var yandexResp struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&yandexResp); err != nil {
		fmt.Printf("[YANDEX ERROR] Ошибка декодирования JSON: %v\n", err)
		return s.getFallbackResult()
	}

	if len(yandexResp.Result.Alternatives) == 0 {
		fmt.Println("[YANDEX ERROR] Яндекс вернул пустой список ответов")
		return s.getFallbackResult()
	}

	rawJSON := yandexResp.Result.Alternatives[0].Message.Text

	// Очищаем от возможных кавычек ```json
	rawJSON = strings.TrimSpace(rawJSON)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimPrefix(rawJSON, "```")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	var result domain.AIResult
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		fmt.Printf("[YANDEX ERROR] Ошибка парсинга JSON структуры ИИ: %v. Текст: %s\n", err, rawJSON)
		return s.getFallbackResult()
	}

	return result
}

func (s *AIService) getFallbackResult() domain.AIResult {
	return domain.AIResult{
		Sentiment: "neutral",
		Reply:     "Здравствуйте! Спасибо за ваше обращение. Я получил ваше сообщение и свяжусь с вами в ближайшее время для обсуждения деталей.",
	}
}
