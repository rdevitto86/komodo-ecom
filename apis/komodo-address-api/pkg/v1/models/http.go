package models

type ValidateRequest struct {
	Address Address `json:"address"`
}

type ValidateResponse struct {
	Valid  bool                `json:"valid"`
	Errors map[string][]string `json:"errors,omitempty"`
}

type NormalizeRequest struct {
	Address Address `json:"address"`
}

type NormalizeResponse struct {
	Address Address `json:"address"`
}

type GeocodeRequest struct {
	Address Address `json:"address"`
}

type GeocodeResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
