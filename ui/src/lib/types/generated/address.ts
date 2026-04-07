// Auto-generated from apis/komodo-address-api/docs/openapi.yaml
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
    "/addresses/validate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Validate an address */
        post: operations["validateAddress"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/addresses/normalize": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Normalize an address to a standardized format */
        post: operations["normalizeAddress"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/addresses/geocode": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Geocode an address to lat/lng coordinates */
        post: operations["geocodeAddress"];
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
        AddressInput: {
            street1: string;
            street2?: string;
            city: string;
            /** @description State/province/region code (e.g. "CA") */
            state: string;
            postalCode: string;
            /**
             * @description ISO country code — defaults to "US" if omitted
             * @default US
             */
            country: string;
        };
        ValidationResult: {
            valid: boolean;
            deliverable: boolean;
            /** @description Field-level validation errors */
            errors?: {
                [key: string]: string;
            };
        };
        NormalizedAddress: {
            street1: string;
            street2?: string;
            city: string;
            state: string;
            postalCode: string;
            country: string;
        };
        NormalizeResponse: {
            address: components["schemas"]["NormalizedAddress"];
        };
        GeocodeResponse: {
            /** Format: double */
            latitude: number;
            /** Format: double */
            longitude: number;
            /** @description Provider accuracy level (e.g. "rooftop", "approximate") */
            accuracy: string;
            /** @description Which geocoding provider resolved the address */
            provider: string;
            normalized?: components["schemas"]["NormalizedAddress"];
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
        /** @description Invalid request body */
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
    validateAddress: {
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
            /** @description Validation result */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ValidationResult"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Address provider unavailable */
            502: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problem"];
                };
            };
        };
    };
    normalizeAddress: {
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
            /** @description Normalized address */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["NormalizeResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Address provider unavailable */
            502: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problem"];
                };
            };
        };
    };
    geocodeAddress: {
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
            /** @description Geocoded result */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["GeocodeResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Address provider unavailable */
            502: {
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
export type AddressInput = components['schemas']['AddressInput'];
export type ValidationResult = components['schemas']['ValidationResult'];
export type NormalizedAddress = components['schemas']['NormalizedAddress'];
export type NormalizeResponse = components['schemas']['NormalizeResponse'];
export type GeocodeResponse = components['schemas']['GeocodeResponse'];
export type Problem = components['schemas']['Problem'];