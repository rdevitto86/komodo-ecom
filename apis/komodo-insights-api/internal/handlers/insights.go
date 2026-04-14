package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-insights-api/internal/service"
)

var svc *service.InsightsService

// InitService wires the InsightsService into the handler layer.
// Must be called from main after the LLM provider is initialised.
func InitService(s *service.InsightsService) {
	svc = s
}

// GetItemSummary handles GET /insights/items/{itemId}/summary.
func GetItemSummary(wtr http.ResponseWriter, req *http.Request) {
	itemID := req.PathValue("itemId")
	if itemID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	result, err := svc.GetItemSummary(req.Context(), itemID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpErr.SendError(wtr, req, httpErr.Global.NotFound)
			return
		}
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(result)
}

// GetItemSentiment handles GET /insights/items/{itemId}/sentiment.
func GetItemSentiment(wtr http.ResponseWriter, req *http.Request) {
	itemID := req.PathValue("itemId")
	if itemID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	result, err := svc.GetItemSentiment(req.Context(), itemID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httpErr.SendError(wtr, req, httpErr.Global.NotFound)
			return
		}
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(result)
}

// GetTrending handles GET /insights/trending.
func GetTrending(wtr http.ResponseWriter, req *http.Request) {
	result, err := svc.GetTrending(req.Context())
	if err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(result)
}
