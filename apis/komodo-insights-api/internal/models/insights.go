package models

import "time"

// SummaryRequest is passed to the LLM provider regardless of backend.
// Context holds the raw text corpus (e.g. review bodies, product descriptions)
// that the provider should summarise.
type SummaryRequest struct {
	EntityType string   // "item" | "service" | "category"
	EntityID   string
	Context    []string // ordered, most-relevant first
}

// ItemSummaryResponse is returned by GET /insights/items/{itemId}/summary.
type ItemSummaryResponse struct {
	ItemID      string    `json:"item_id"`
	Summary     string    `json:"summary"`
	GeneratedAt time.Time `json:"generated_at"`
}

// SentimentBreakdown holds normalised sentiment percentages (0–100, must sum to 100).
type SentimentBreakdown struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}

// ItemSentimentResponse is returned by GET /insights/items/{itemId}/sentiment.
type ItemSentimentResponse struct {
	ItemID      string             `json:"item_id"`
	Sentiment   SentimentBreakdown `json:"sentiment"`
	TopThemes   []string           `json:"top_themes"`
	ReviewCount int                `json:"review_count"`
	GeneratedAt time.Time          `json:"generated_at"`
}

// TrendingItem represents a single item in the trending list.
type TrendingItem struct {
	ItemID string  `json:"item_id"`
	Signal string  `json:"signal"` // e.g. "rising_reviews", "high_rating", "bestseller"
	Score  float64 `json:"score"`
}

// TrendingResponse is returned by GET /insights/trending.
type TrendingResponse struct {
	Items       []TrendingItem `json:"items"`
	GeneratedAt time.Time      `json:"generated_at"`
}
