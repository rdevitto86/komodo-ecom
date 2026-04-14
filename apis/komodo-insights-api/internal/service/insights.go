package service

import (
	"context"
	"fmt"

	"komodo-insights-api/internal/models"
)

// InsightsService orchestrates data fetching and LLM summarisation.
// At stub stage all methods return ErrNotFound — wire the provider and
// upstream clients (reviews-api, shop-items-api) before implementing.
type InsightsService struct {
	provider SummaryProvider
}

func NewInsightsService(provider SummaryProvider) *InsightsService {
	return &InsightsService{provider: provider}
}

// GetItemSummary fetches item reviews + product details, then asks the
// provider for a concise natural-language summary.
//
// TODO: call reviews-api GET /items/{itemId}/reviews (paginated, top N by recency + rating)
// TODO: call shop-items-api GET /items/{itemId} for product description
// TODO: build SummaryRequest from corpus and call provider.Summarize
func (s *InsightsService) GetItemSummary(ctx context.Context, itemID string) (*models.ItemSummaryResponse, error) {
	return nil, fmt.Errorf("GetItemSummary: %w", ErrNotFound)
}

// GetItemSentiment derives a sentiment breakdown and top themes from item reviews.
//
// TODO: fetch reviews from reviews-api; pass to provider for sentiment extraction
// TODO: normalise Positive/Neutral/Negative scores so they sum to 100
func (s *InsightsService) GetItemSentiment(ctx context.Context, itemID string) (*models.ItemSentimentResponse, error) {
	return nil, fmt.Errorf("GetItemSentiment: %w", ErrNotFound)
}

// GetTrending returns a ranked list of trending items across the catalog.
//
// TODO: define signal sources — candidates: review velocity (reviews-api),
//       order frequency (order-api), add-to-cart rate (cart-api)
// TODO: decide whether trending is LLM-derived or signal-aggregated (or both)
func (s *InsightsService) GetTrending(ctx context.Context) (*models.TrendingResponse, error) {
	return nil, fmt.Errorf("GetTrending: %w", ErrNotFound)
}
