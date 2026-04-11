package models

// SearchResult represents a single item returned from the Typesense index.
// Products and services share the same result shape — use Type to differentiate.
//
// TODO(typesense): align field names with the Typesense collection schema
// defined in repository.CollectionSchema.
type SearchResult struct {
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`     // "product" | "service"
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	ImageURL    string   `json:"image_url,omitempty"`
	Score       float64  `json:"score"` // Typesense relevance score
}

// SearchResponse is the paginated response envelope for GET /search.
type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}
