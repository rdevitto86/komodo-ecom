package subscriber

import (
	"context"

	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// StartShopItemsSubscriber starts a background subscriber that listens for shop item
// events from events-api and syncs changes to the Typesense index.
//
// Called from main() after Typesense client is initialized.
//
// TODO(events): wire up to events-api pub/sub transport once events-api is ready.
// The transport (Redis Streams, SNS/SQS, or custom) is TBD — confirm with events-api design.
//
// Expected event types to handle:
//   - shop_item.created  → repository.IndexItem(ctx, toSearchResult(event))
//   - shop_item.updated  → repository.IndexItem(ctx, toSearchResult(event))
//   - shop_item.deleted  → repository.DeleteItem(ctx, event.SKU)
//
// TODO(events): implement toSearchResult(event) — maps the shop-item event payload
// to models.SearchResult. Source of truth for field mapping is shop-items-api event schema.
//
// TODO(resilience): on subscriber failure, log error and attempt reconnect with backoff.
// Do not crash the service — search reads should still work if sync is temporarily down.
func StartShopItemsSubscriber(ctx context.Context) {
	logger.Info("shop items subscriber: starting")

	// TODO: connect to events-api transport
	// TODO: subscribe to shop_item.* event types
	// TODO: dispatch to repository.IndexItem / repository.DeleteItem based on event type

	logger.Warn("shop items subscriber: not yet implemented — Typesense index will not sync until this is wired")
}
