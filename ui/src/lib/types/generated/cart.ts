// Auto-generated from apis/komodo-cart-api/docs/openapi.yaml
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
    "/me/cart": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get authenticated user's cart */
        get: operations["getMyCart"];
        put?: never;
        post?: never;
        /** Clear authenticated user's cart */
        delete: operations["clearMyCart"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/cart/merge": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Merge a guest cart into the authenticated cart
         * @description Merges items from a guest cart (Redis) into the authenticated user's cart (DynamoDB).
         *     Quantities are additive for duplicate items; authenticated cart wins on price/name conflicts.
         *     The guest cart Redis key is deleted after a successful merge.
         *     Returns the merged cart — no separate GET needed.
         */
        post: operations["mergeGuestCart"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/cart/items": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Add item to authenticated cart */
        post: operations["addMyCartItem"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/cart/items/{itemId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        /** Update item quantity in authenticated cart */
        put: operations["updateMyCartItem"];
        post?: never;
        /** Remove item from authenticated cart */
        delete: operations["deleteMyCartItem"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/cart/checkout": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Initiate checkout — places stock holds, returns checkout token */
        post: operations["initiateCheckout"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/cart": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Create a guest cart */
        post: operations["createGuestCart"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/cart/{cartId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get guest cart */
        get: operations["getGuestCart"];
        put?: never;
        post?: never;
        /** Clear guest cart */
        delete: operations["clearGuestCart"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/cart/{cartId}/items": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Add item to guest cart */
        post: operations["addGuestCartItem"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/cart/{cartId}/items/{itemId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        /** Update item quantity in guest cart */
        put: operations["updateGuestCartItem"];
        post?: never;
        /** Remove item from guest cart */
        delete: operations["deleteGuestCartItem"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        Cart: {
            /** Format: uuid */
            id: string;
            /** @description Present for authenticated carts only */
            user_id?: string;
            items: components["schemas"]["CartItem"][];
            /** @description Sum of (unit_price_cents * quantity) for all items */
            subtotal_cents: number;
            item_count: number;
            /** Format: date-time */
            updated_at: string;
        };
        CartItem: {
            /** @description Product ID from shop-items-api */
            item_id: string;
            /** @description Variant SKU */
            sku: string;
            /** @description Product name snapshot at add-time */
            name: string;
            quantity: number;
            /** @description Price snapshot at add-time */
            unit_price_cents: number;
            image_url?: string;
        };
        MergeCartRequest: {
            /**
             * Format: uuid
             * @description UUID of the guest cart to merge. Generated client-side by the UI,
             *     stored in localStorage, and sent here after login.
             */
            guest_cart_id: string;
        };
        AddItemRequest: {
            item_id: string;
            sku: string;
            quantity: number;
        };
        UpdateItemRequest: {
            /** @description Set to 0 to remove the item */
            quantity: number;
        };
        CheckoutResponse: {
            /** @description Short-lived token passed to order-api to confirm purchase */
            checkout_token: string;
            /**
             * Format: date-time
             * @description When the stock holds expire (CART_HOLD_TTL_SEC from now)
             */
            expires_at: string;
            items: components["schemas"]["CartItem"][];
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
        /** @description Session token does not match cart */
        Forbidden: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["Problem"];
            };
        };
        /** @description Cart or item not found */
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
        CartId: string;
        ItemId: string;
        /** @description Guest session token (set by client on cart creation, stored in cookie) */
        SessionId: string;
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
    getMyCart: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Cart retrieved */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    clearMyCart: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Cart cleared */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    mergeGuestCart: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["MergeCartRequest"];
            };
        };
        responses: {
            /** @description Guest cart merged, merged cart returned */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Guest cart not found (missing or expired — both return 404) */
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
    addMyCartItem: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["AddItemRequest"];
            };
        };
        responses: {
            /** @description Item added */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
    updateMyCartItem: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                itemId: components["parameters"]["ItemId"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateItemRequest"];
            };
        };
        responses: {
            /** @description Item updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    deleteMyCartItem: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                itemId: components["parameters"]["ItemId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Item removed */
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
    initiateCheckout: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Checkout initiated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CheckoutResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Insufficient stock for one or more items */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Problem"];
                };
            };
        };
    };
    createGuestCart: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Guest cart created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
        };
    };
    getGuestCart: {
        parameters: {
            query?: never;
            header: {
                /** @description Guest session token (set by client on cart creation, stored in cookie) */
                "X-Session-ID": components["parameters"]["SessionId"];
            };
            path: {
                cartId: components["parameters"]["CartId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Cart retrieved */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            403: components["responses"]["Forbidden"];
            404: components["responses"]["NotFound"];
        };
    };
    clearGuestCart: {
        parameters: {
            query?: never;
            header: {
                /** @description Guest session token (set by client on cart creation, stored in cookie) */
                "X-Session-ID": components["parameters"]["SessionId"];
            };
            path: {
                cartId: components["parameters"]["CartId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Cart cleared */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            403: components["responses"]["Forbidden"];
            404: components["responses"]["NotFound"];
        };
    };
    addGuestCartItem: {
        parameters: {
            query?: never;
            header: {
                /** @description Guest session token (set by client on cart creation, stored in cookie) */
                "X-Session-ID": components["parameters"]["SessionId"];
            };
            path: {
                cartId: components["parameters"]["CartId"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["AddItemRequest"];
            };
        };
        responses: {
            /** @description Item added */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            400: components["responses"]["BadRequest"];
            403: components["responses"]["Forbidden"];
            404: components["responses"]["NotFound"];
        };
    };
    updateGuestCartItem: {
        parameters: {
            query?: never;
            header: {
                /** @description Guest session token (set by client on cart creation, stored in cookie) */
                "X-Session-ID": components["parameters"]["SessionId"];
            };
            path: {
                cartId: components["parameters"]["CartId"];
                itemId: components["parameters"]["ItemId"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateItemRequest"];
            };
        };
        responses: {
            /** @description Item updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Cart"];
                };
            };
            400: components["responses"]["BadRequest"];
            403: components["responses"]["Forbidden"];
            404: components["responses"]["NotFound"];
        };
    };
    deleteGuestCartItem: {
        parameters: {
            query?: never;
            header: {
                /** @description Guest session token (set by client on cart creation, stored in cookie) */
                "X-Session-ID": components["parameters"]["SessionId"];
            };
            path: {
                cartId: components["parameters"]["CartId"];
                itemId: components["parameters"]["ItemId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Item removed */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            403: components["responses"]["Forbidden"];
            404: components["responses"]["NotFound"];
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type Cart = components['schemas']['Cart'];
export type CartItem = components['schemas']['CartItem'];
export type MergeCartRequest = components['schemas']['MergeCartRequest'];
export type AddItemRequest = components['schemas']['AddItemRequest'];
export type UpdateItemRequest = components['schemas']['UpdateItemRequest'];
export type CheckoutResponse = components['schemas']['CheckoutResponse'];
export type Problem = components['schemas']['Problem'];