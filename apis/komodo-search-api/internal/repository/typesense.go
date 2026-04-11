package repository

import (
	"context"

	"komodo-search-api/internal/models"
)

// TODO(typesense): add dependency to go.mod:
//   go get github.com/typesense/typesense-go
//
// CollectionSchema defines the Typesense collection for shop items (products + services).
// Run this schema once on startup (or via POST /internal/index/sync) to create the collection.
//
//	var CollectionSchema = &typesense.CollectionSchema{
//	    Name: collectionName,
//	    Fields: []typesense.Field{
//	        {Name: "sku",         Type: "string",  Facet: false},
//	        {Name: "name",        Type: "string",  Facet: false},
//	        {Name: "description", Type: "string",  Facet: false},
//	        {Name: "type",        Type: "string",  Facet: true},  // "product" | "service"
//	        {Name: "price",       Type: "float",   Facet: true},
//	        {Name: "category",    Type: "string",  Facet: true},
//	        {Name: "tags",        Type: "string[]", Facet: true},
//	        {Name: "is_active",   Type: "bool",    Facet: false},
//	    },
//	    DefaultSortingField: "name",
//	}

// TODO(typesense): package-level Typesense client, initialized by InitTypesense.
// var tsClient *typesense.Client
// var collectionName string

// InitTypesense initializes the Typesense client and verifies the collection exists.
//
// TODO(typesense): implement with typesense-go client:
//   tsClient = typesense.NewClient(
//       typesense.WithServer(fmt.Sprintf("http://%s:%s", host, port)),
//       typesense.WithAPIKey(apiKey),
//   )
//   collectionName = collection
//   verify collection exists, create if missing using CollectionSchema
func InitTypesense(host, port, apiKey, collection string) error {
	// TODO: implement
	return nil
}

// SearchParams defines the query parameters for a search request.
type SearchParams struct {
	Query    string
	ItemType string // "product" | "service" | "" (all)
	Category string
	Page     int
	PerPage  int
}

// Search executes a search query against the Typesense index.
//
// TODO(typesense): implement using typesense-go SearchParameters:
//
//	searchParams := &typesense.SearchCollectionParams{
//	    Q:        params.Query,
//	    QueryBy:  "name,description,tags",
//	    Page:     params.Page,
//	    PerPage:  params.PerPage,
//	}
//	if params.ItemType != "" {
//	    searchParams.FilterBy = "type:=" + params.ItemType
//	}
//	if params.Category != "" {
//	    searchParams.FilterBy += " && category:=" + params.Category
//	}
//
// Map Typesense hits to []models.SearchResult. Include _text_match score as Score field.
func Search(ctx context.Context, params SearchParams) ([]models.SearchResult, int, error) {
	// TODO: implement Typesense query
	return []models.SearchResult{}, 0, nil
}

// IndexItem adds or updates a single item in the Typesense collection.
// Called by the events-api subscriber when a shop item is created or updated.
//
// TODO(typesense): use tsClient.Collection(collectionName).Documents().Upsert(ctx, doc)
// where doc is built from the shop item event payload.
func IndexItem(ctx context.Context, result models.SearchResult) error {
	// TODO: implement upsert
	return nil
}

// DeleteItem removes an item from the Typesense collection by SKU.
// Called by the events-api subscriber when a shop item is deleted.
//
// TODO(typesense): use tsClient.Collection(collectionName).Document(sku).Delete(ctx)
func DeleteItem(ctx context.Context, sku string) error {
	// TODO: implement delete
	return nil
}
