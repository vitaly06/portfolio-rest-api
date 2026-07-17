package repository

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/vitaly06/portfolio-rest-api/internal/domain"
)

type FileRepository struct {
	mu          sync.Mutex
	logFile     string
	metricsFile string
}

func NewFileRepository(logFile, metricsFile string) *FileRepository {
	return &FileRepository{
		logFile:     logFile,
		metricsFile: metricsFile,
	}
}

func (r *FileRepository) WriteLog(entry domain.LogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(data) + "\n")
	return err
}

func (r *FileRepository) UpdateMetrics(sentiment string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var metrics domain.Metrics

	// Читаем старые метрики, если файл существует
	if data, err := os.ReadFile(r.metricsFile); err == nil {
		_ = json.Unmarshal(data, &metrics)
	}

	metrics.TotalRequests++

	// Инкрементируем нужное поле в зависимости от тональности
	switch sentiment {
	case "positive":
		metrics.SentimentStats.Positive++
	case "negative":
		metrics.SentimentStats.Negative++
	default:
		metrics.SentimentStats.Neutral++
	}

	metrics.LastUpdate = time.Now()

	updatedData, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.metricsFile, updatedData, 0644)
}
func (r *FileRepository) GetMetrics() (domain.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Инициализируем пустую структуру (поля интов автоматически будут равны 0)
	var metrics domain.Metrics

	data, err := os.ReadFile(r.metricsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return metrics, nil
		}
		return metrics, err
	}

	err = json.Unmarshal(data, &metrics)
	return metrics, err
}
