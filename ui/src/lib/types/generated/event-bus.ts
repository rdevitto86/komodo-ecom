// Auto-generated from apis/komodo-event-bus-api/docs/openapi.yaml
// Do not edit — run `bun run gen:types` to regenerate.

export interface paths {
    "/health": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Liveness check */
        get: operations["getHealth"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/events": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Publish an event */
        post: operations["publishEvent"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        EventEnvelope: {
            /**
             * Format: uuid
             * @description Idempotency key — consumers must deduplicate on this
             */
            id: string;
            /**
             * @description <domain>.<action> e.g. order.placed
             * @example order.placed
             */
            type: string;
            /**
             * @description Numeric string — increment when payload shape changes
             * @example 1
             */
            version: string;
            /**
             * @description Publishing service name
             * @example komodo-order-api
             */
            source: string;
            /**
             * Format: date-time
             * @description When the event happened (not when it was published)
             */
            occurred_at: string;
            /** @description Event-specific data — typed per event type and version */
            payload: Record<string, never>;
        };
        PublishResponse: {
            /**
             * Format: uuid
             * @description Echo of the submitted event ID
             */
            event_id: string;
            /** @description SNS message ID (or local correlation ID in dev) */
            message_id: string;
        };
        /** @description RFC 7807 Problem+JSON */
        Problem: {
            /** Format: uri */
            type: string;
            title: string;
            status: number;
            detail?: string;
            instance?: string;
        };
    };
    responses: {
        /** @description Invalid request body or parameters */
        BadRequest: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["Problem"];
            };
        };
        /** @description Missing or invalid JWT */
        Unauthorized: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["Problem"];
            };
        };
    };
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export interface operations {
    getHealth: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    publishEvent: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EventEnvelope"];
            };
        };
        responses: {
            /** @description Event accepted for delivery */
            202: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PublishResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Unknown event type or invalid payload for declared version */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problem"];
                };
            };
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type EventEnvelope = components['schemas']['EventEnvelope'];
export type PublishResponse = components['schemas']['PublishResponse'];
export type Problem = components['schemas']['Problem'];