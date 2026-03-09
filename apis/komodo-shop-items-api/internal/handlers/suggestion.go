package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
	shopitems "komodo-forge-sdk-go/http/services/shop_items"
	logger "komodo-forge-sdk-go/logging/runtime"
)

// Returns product suggestions based on user viewing habits (auth required)
func GetSuggestions(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	var reqBody shopitems.SuggestionRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse suggestion request body", err)
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to parse request body"))
		return
	}

	if reqBody.Limit <= 0 {
		reqBody.Limit = 10
	}

	// TODO: integrate with recommendation engine / ML service
	// For now, return empty suggestions
	response := shopitems.SuggestionResponse{
		Suggestions: []shopitems.Product{},
		Total:       0,
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(response)
}
