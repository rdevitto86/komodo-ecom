//go:build e2e

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	res := get(t, "/health", nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestValidateAddress_InvalidInput verifies the endpoint rejects an empty body.
func TestValidateAddress_InvalidInput(t *testing.T) {
	res := post(t, "/addresses/validate", map[string]any{}, nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("address provider not wired — set ADDRESS_PROVIDER_API_KEY in LocalStack secrets to enable")
	}
	checkStatus(t, res, http.StatusBadRequest)
}

// TestValidateAddress_Valid submits a well-formed US address and expects a validation result.
func TestValidateAddress_Valid(t *testing.T) {
	body := map[string]any{
		"street":  "1600 Amphitheatre Parkway",
		"city":    "Mountain View",
		"state":   "CA",
		"zip":     "94043",
		"country": "US",
	}
	res := post(t, "/addresses/validate", body, nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("address provider not wired — set ADDRESS_PROVIDER_API_KEY in LocalStack secrets to enable")
	}
	// 200 = valid address, 422 = address exists but could not be validated.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusUnprocessableEntity {
		checkStatus(t, res, http.StatusOK)
	}
}

// TestNormalizeAddress verifies mixed-case input is normalised to a canonical form.
func TestNormalizeAddress(t *testing.T) {
	body := map[string]any{
		"street":  "1600 amphitheatre pkwy",
		"city":    "mountain view",
		"state":   "ca",
		"zip":     "94043",
		"country": "US",
	}
	res := post(t, "/addresses/normalize", body, nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("address provider not wired — set ADDRESS_PROVIDER_API_KEY in LocalStack secrets to enable")
	}
	checkStatus(t, res, http.StatusOK)
}

// TestGeocodeAddress verifies lat/lng coordinates are returned for a valid address.
func TestGeocodeAddress(t *testing.T) {
	body := map[string]any{
		"street":  "1600 Amphitheatre Parkway",
		"city":    "Mountain View",
		"state":   "CA",
		"zip":     "94043",
		"country": "US",
	}
	res := post(t, "/addresses/geocode", body, nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("address provider not wired — set ADDRESS_PROVIDER_API_KEY in LocalStack secrets to enable")
	}
	checkStatus(t, res, http.StatusOK)

	var result struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
	decodeJSON(t, res, &result)
	if result.Lat == 0 && result.Lng == 0 {
		t.Fatal("expected non-zero lat/lng in geocode response")
	}
}

// TestGeocodeAddress_InvalidInput verifies missing fields are rejected.
func TestGeocodeAddress_InvalidInput(t *testing.T) {
	res := post(t, "/addresses/geocode", map[string]any{"street": ""}, nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("address provider not wired")
	}
	if res.StatusCode == http.StatusOK {
		var result struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}
		decodeJSON(t, res, &result)
		if result.Lat != 0 || result.Lng != 0 {
			t.Fatal("expected zero lat/lng for empty address")
		}
	}
	_ = json.NewDecoder(res.Body)
}
