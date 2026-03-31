package provider

import (
	"context"

	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

type AddressInput struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country,omitempty"`
}

type ValidationResult struct {
	Valid        bool              `json:"valid"`
	Deliverable  bool              `json:"deliverable"`
	Errors       map[string]string `json:"errors,omitempty"`
}

type NormalizedAddress struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type GeocodedAddress struct {
	Latitude   float64            `json:"latitude"`
	Longitude  float64            `json:"longitude"`
	Accuracy   string             `json:"accuracy"`
	Provider   string             `json:"provider"`
	Normalized *NormalizedAddress `json:"normalized,omitempty"`
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// TODO: replace stub body with real provider HTTP call (e.g. SmartyStreets or Google Maps validation endpoint)
func (c *Client) Validate(ctx context.Context, addr AddressInput) (*ValidationResult, error) {
	logger.Warn("address provider not wired — returning stub validation result", logger.FromContext(ctx)...)
	return &ValidationResult{
		Valid:       true,
		Deliverable: true,
	}, nil
}

// TODO: replace stub body with real provider HTTP call for address standardization
func (c *Client) Normalize(ctx context.Context, addr AddressInput) (*NormalizedAddress, error) {
	logger.Warn("address provider not wired — returning stub normalization result", logger.FromContext(ctx)...)

	country := addr.Country
	if country == "" {
		country = "US"
	}
	return &NormalizedAddress{
		Street1:    addr.Street1,
		Street2:    addr.Street2,
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.PostalCode,
		Country:    country,
	}, nil
}

// TODO: replace stub body with real provider HTTP call for geocoding (e.g. Google Maps Geocoding API)
func (c *Client) Geocode(ctx context.Context, addr AddressInput) (*GeocodedAddress, error) {
	logger.Warn("address provider not wired — returning stub geocode result", logger.FromContext(ctx)...)

	country := addr.Country
	if country == "" {
		country = "US"
	}
	return &GeocodedAddress{
		Latitude:  0.0,
		Longitude: 0.0,
		Accuracy:  "stub",
		Provider:  "stub",
		Normalized: &NormalizedAddress{
			Street1:    addr.Street1,
			Street2:    addr.Street2,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    country,
		},
	}, nil
}
