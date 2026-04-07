// Auto-generated from apis/komodo-shop-inventory-api/docs/openapi.yaml
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
    "/stock": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Batch stock levels for a list of SKUs */
        get: operations["getBatchStock"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/stock/{sku}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get stock level for a single SKU */
        get: operations["getStock"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/stock/{sku}/reserve": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Place a stock hold */
        post: operations["reserveStock"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/stock/{sku}/reserve/{holdId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /** Release a stock hold early */
        delete: operations["releaseHold"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/stock/{sku}/confirm": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Convert hold to confirmed decrement on order confirmation */
        post: operations["confirmStock"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/stock/{sku}/restock": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Increment stock (new stock received or return processed) */
        post: operations["restockSku"];
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
        StockLevel: {
            sku: string;
            /** @description Units available for new reservations */
            available_qty: number;
            /** @description Units currently held (not yet confirmed) */
            reserved_qty: number;
            /** @description Units sold (orders confirmed) */
            committed_qty: number;
            /** @description Alert fires when available_qty drops below this */
            restock_threshold?: number;
        };
        ReserveRequest: {
            /** Format: uuid */
            cart_id: string;
            quantity: number;
        };
        HoldResponse: {
            /** Format: uuid */
            hold_id: string;
            sku: string;
            quantity: number;
            /** Format: date-time */
            expires_at: string;
        };
        ConfirmRequest: {
            /** Format: uuid */
            hold_id: string;
        };
        RestockRequest: {
            quantity: number;
            /**
             * @description e.g. "purchase_received", "return_processed"
             * @enum {string}
             */
            reason?: "purchase_received" | "return_processed" | "adjustment";
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
        /** @description SKU or hold not found */
        NotFound: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["Problem"];
            };
        };
    };
    parameters: {
        Sku: string;
        HoldId: string;
    };
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
    getBatchStock: {
        parameters: {
            query: {
                /** @description Comma-separated SKUs (max 100) */
                skus: string[];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Stock levels */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": {
                        [key: string]: components["schemas"]["StockLevel"];
                    };
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
    getStock: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                sku: components["parameters"]["Sku"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Stock level */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["StockLevel"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    reserveStock: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                sku: components["parameters"]["Sku"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ReserveRequest"];
            };
        };
        responses: {
            /** @description Hold placed */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HoldResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Insufficient stock */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problem"];
                };
            };
        };
    };
    releaseHold: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                sku: components["parameters"]["Sku"];
                holdId: components["parameters"]["HoldId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Hold released */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    confirmStock: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                sku: components["parameters"]["Sku"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ConfirmRequest"];
            };
        };
        responses: {
            /** @description Stock confirmed */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Hold not found or already expired */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problem"];
                };
            };
        };
    };
    restockSku: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                sku: components["parameters"]["Sku"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["RestockRequest"];
            };
        };
        responses: {
            /** @description Stock updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["StockLevel"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type StockLevel = components['schemas']['StockLevel'];
export type ReserveRequest = components['schemas']['ReserveRequest'];
export type HoldResponse = components['schemas']['HoldResponse'];
export type ConfirmRequest = components['schemas']['ConfirmRequest'];
export type RestockRequest = components['schemas']['RestockRequest'];
export type Problem = components['schemas']['Problem'];