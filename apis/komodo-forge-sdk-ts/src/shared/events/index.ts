// EventType is the canonical event name in <entity>.<verb> format.
// This is the public contract between services — treat as stable once defined.
export type EventType =
  | 'order.created'
  | 'order.status_updated'
  | 'order.cancelled'
  | 'order.fulfilled'
  | 'user.created'
  | 'user.profile_updated'
  | 'user.deleted'
  | 'payment.initiated'
  | 'payment.succeeded'
  | 'payment.failed'
  | 'payment.refunded'
  | 'cart.checked_out'
  | 'inventory.reserved'
  | 'inventory.released';

// EventSource identifies which Komodo service emitted the event.
export type EventSource =
  | 'komodo-auth-api'
  | 'komodo-user-api'
  | 'komodo-order-api'
  | 'komodo-cart-api'
  | 'komodo-inventory-api'
  | 'komodo-payments-api'
  | 'komodo-shop-items-api'
  | 'komodo-communications-api';

// EntityType identifies the domain entity at the centre of the event.
export type EntityType =
  | 'order'
  | 'user'
  | 'payment'
  | 'cart'
  | 'inventory'
  | 'product';

// KomodoEvent is the canonical business event envelope published to SNS FIFO
// topics and consumed via SQS FIFO queues. The SNS/SQS MessageGroupId must be
// set to entity_id so that all events for the same entity are ordered, while
// events for different entities can be processed in parallel.
//
// Field names match the wire format (snake_case) defined in the event-bus-api
// OpenAPI spec. The generic T parameter narrows the payload for known event types.
export interface KomodoEvent<T = Record<string, unknown>> {
  id: string;
  type: EventType;
  source: EventSource;
  entity_id: string;
  entity_type: EntityType;
  /** ISO 8601 UTC timestamp */
  occurred_at: string;
  version: string;
  payload: T;
  /** Traces the event chain back to the originating HTTP request. */
  correlation_id?: string;
}
