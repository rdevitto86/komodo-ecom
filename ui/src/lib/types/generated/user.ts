// Auto-generated from apis/komodo-user-api/docs/openapi.yaml
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
        get: operations["healthCheck"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/profile": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get authenticated user's profile */
        get: operations["getProfile"];
        /** Update authenticated user's profile */
        put: operations["updateProfile"];
        /** Create user record (called on registration) */
        post: operations["createUser"];
        /** Delete authenticated user's account */
        delete: operations["deleteProfile"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/addresses": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List all addresses for the authenticated user */
        get: operations["getAddresses"];
        put?: never;
        /** Add a new address */
        post: operations["addAddress"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/addresses/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Address ID */
                id: components["parameters"]["AddressID"];
            };
            cookie?: never;
        };
        get?: never;
        /** Update an address */
        put: operations["updateAddress"];
        post?: never;
        /** Delete an address */
        delete: operations["deleteAddress"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/payments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List saved payment methods */
        get: operations["getPayments"];
        /** Add or update a payment method */
        put: operations["upsertPayment"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/payments/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Payment method ID */
                id: components["parameters"]["PaymentID"];
            };
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /** Remove a payment method */
        delete: operations["deletePayment"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/preferences": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get user preferences */
        get: operations["getPreferences"];
        /** Update user preferences */
        put: operations["updatePreferences"];
        post?: never;
        /** Delete user preferences */
        delete: operations["deletePreferences"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        /** Get profile by user ID (service-to-service) */
        get: operations["internalGetUser"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users/{id}/addresses": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        /** Get addresses for a user (service-to-service) */
        get: operations["internalGetUserAddresses"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users/{id}/preferences": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        /** Get preferences for a user (service-to-service) */
        get: operations["internalGetUserPreferences"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/users/{id}/payments": {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        /** Get payment methods for a user (service-to-service) */
        get: operations["internalGetUserPayments"];
        put?: never;
        post?: never;
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
        User: {
            user_id: string;
            /** Format: email */
            email: string;
            phone?: string;
            first_name: string;
            middle_initial?: string;
            last_name: string;
            /** Format: uri */
            avatar_url?: string;
            /** Format: date-time */
            readonly created_at: string;
            /** Format: date-time */
            readonly updated_at: string;
        };
        CreateUserRequest: {
            /** Format: email */
            email: string;
            phone?: string;
            first_name: string;
            middle_initial?: string;
            last_name: string;
            /** Format: uri */
            avatar_url?: string;
        };
        UpdateUserRequest: {
            phone?: string;
            first_name?: string;
            middle_initial?: string;
            last_name?: string;
            /** Format: uri */
            avatar_url?: string;
        };
        Address: {
            address_id: string;
            /** @description Friendly label (e.g. "Home", "Work") */
            alias?: string;
            line1: string;
            line2?: string;
            city: string;
            state: string;
            zip_code: string;
            /** @description ISO 3166-1 alpha-2 country code (e.g. "US") */
            country: string;
            is_default: boolean;
        };
        AddressInput: {
            alias?: string;
            line1: string;
            line2?: string;
            city: string;
            state: string;
            zip_code: string;
            /** @description ISO 3166-1 alpha-2 country code (e.g. "US") */
            country: string;
            is_default?: boolean;
        };
        Preferences: {
            /** @description BCP 47 language tag (e.g. "en-US") */
            language?: string;
            /** @description IANA timezone (e.g. "America/New_York") */
            timezone?: string;
            /** @description Notification channel opt-ins keyed by channel name (e.g. email, sms, push) */
            communication?: {
                [key: string]: boolean;
            };
            /** @description Marketing preferences keyed by attribute (e.g. frequency, channel) */
            marketing?: {
                [key: string]: string;
            };
        };
        /** @description Stored payment method reference — no raw card data, metadata only */
        PaymentMethod: {
            payment_id: string;
            /** @description Payment processor (e.g. "stripe") */
            provider: string;
            /** @description Last 4 digits of the card number */
            last4: string;
            /** @description Card network (e.g. "visa", "mastercard", "amex") */
            brand: string;
            expiry_month: number;
            expiry_year: number;
            is_default: boolean;
        };
        /** @description Payment processor token submitted by the client to save a payment method */
        PaymentMethodInput: {
            /** @description Payment processor (e.g. "stripe") */
            provider: string;
            /** @description Payment processor token (e.g. Stripe payment method ID "pm_xxx") */
            token: string;
            is_default?: boolean;
        };
        /** @description RFC 7807 Problem+JSON error body */
        ProblemDetail: {
            /** Format: uri */
            type?: string;
            title?: string;
            status?: number;
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
                "application/problem+json": components["schemas"]["ProblemDetail"];
            };
        };
        /** @description Missing or invalid JWT */
        Unauthorized: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["ProblemDetail"];
            };
        };
        /** @description JWT present but insufficient scope */
        Forbidden: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["ProblemDetail"];
            };
        };
        /** @description Resource not found */
        NotFound: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["ProblemDetail"];
            };
        };
        /** @description Resource already exists */
        Conflict: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["ProblemDetail"];
            };
        };
    };
    parameters: {
        /** @description User ID */
        UserID: string;
        /** @description Address ID */
        AddressID: string;
        /** @description Payment method ID */
        PaymentID: string;
    };
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export interface operations {
    healthCheck: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Service is healthy */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getProfile: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description User profile */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["User"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    updateProfile: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateUserRequest"];
            };
        };
        responses: {
            /** @description Updated user profile */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["User"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    createUser: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateUserRequest"];
            };
        };
        responses: {
            /** @description User created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["User"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            409: components["responses"]["Conflict"];
        };
    };
    deleteProfile: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Account deleted */
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
    getAddresses: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Address list */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Address"][];
                };
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    addAddress: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["AddressInput"];
            };
        };
        responses: {
            /** @description Address created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Address"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
    updateAddress: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Address ID */
                id: components["parameters"]["AddressID"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["AddressInput"];
            };
        };
        responses: {
            /** @description Updated address */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Address"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    deleteAddress: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Address ID */
                id: components["parameters"]["AddressID"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Address deleted */
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
    getPayments: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Payment method list */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PaymentMethod"][];
                };
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    upsertPayment: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PaymentMethodInput"];
            };
        };
        responses: {
            /** @description Payment method saved */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PaymentMethod"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
    deletePayment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Payment method ID */
                id: components["parameters"]["PaymentID"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Payment method removed */
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
    getPreferences: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description User preferences */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Preferences"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    updatePreferences: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["Preferences"];
            };
        };
        responses: {
            /** @description Updated preferences */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Preferences"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
    deletePreferences: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Preferences deleted */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    internalGetUser: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description User profile */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["User"];
                };
            };
            401: components["responses"]["Unauthorized"];
            403: components["responses"]["Forbidden"];
            404: components["responses"]["NotFound"];
        };
    };
    internalGetUserAddresses: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Address list */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Address"][];
                };
            };
            401: components["responses"]["Unauthorized"];
            403: components["responses"]["Forbidden"];
        };
    };
    internalGetUserPreferences: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description User preferences */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Preferences"];
                };
            };
            401: components["responses"]["Unauthorized"];
            403: components["responses"]["Forbidden"];
        };
    };
    internalGetUserPayments: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description User ID */
                id: components["parameters"]["UserID"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Payment method list */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PaymentMethod"][];
                };
            };
            401: components["responses"]["Unauthorized"];
            403: components["responses"]["Forbidden"];
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type User = components['schemas']['User'];
export type CreateUserRequest = components['schemas']['CreateUserRequest'];
export type UpdateUserRequest = components['schemas']['UpdateUserRequest'];
export type Address = components['schemas']['Address'];
export type AddressInput = components['schemas']['AddressInput'];
export type Preferences = components['schemas']['Preferences'];
export type PaymentMethod = components['schemas']['PaymentMethod'];
export type PaymentMethodInput = components['schemas']['PaymentMethodInput'];
export type ProblemDetail = components['schemas']['ProblemDetail'];