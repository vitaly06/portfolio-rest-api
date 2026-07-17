package domain

import "time"

type ContactRequest struct {
	Name    string `json:"name" validate:"required,min=2"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"required"`
	Comment string `json:"comment" validate:"required,min=5"`
}

type AIResult struct {
	Sentiment string `json:"sentiment"`
	Reply     string `json:"reply"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	IP        string `json:"ip"`
	Status    int    `json:"status"`
}

type Metrics struct {
	TotalRequests  int            `json:"total_requests"`
	SentimentStats map[string]int `json:"sentiment_stats"`
	LastUpdate     time.Time      `json:"last_update"`
}
