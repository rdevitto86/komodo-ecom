// Auto-generated from apis/komodo-auth-api/docs/openapi.yaml
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
    "/.well-known/jwks.json": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** JSON Web Key Set — public keys for JWT verification */
        get: operations["getJWKS"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/oauth/token": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Issue an access token
         * @description Supports `client_credentials` and `refresh_token` grant types.
         *     `authorization_code` grant type also accepted when `code` and `redirectUri` are present.
         */
        post: operations["issueToken"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/oauth/authorize": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Authorization endpoint (authorization_code flow) */
        get: operations["authorize"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/oauth/introspect": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Introspect a token (requires client credentials) */
        post: operations["introspectToken"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/oauth/revoke": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Revoke a token (requires client credentials) */
        post: operations["revokeToken"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/internal/token/validate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Validate and decode a JWT (service-to-service)
         * @description Used by other services to validate tokens without implementing JWT verification themselves.
         *     No client JWT required — VPC access is the control.
         *     Requiring auth here would create a circular dependency.
         */
        post: operations["validateToken"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/internal/clients": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List all registered OAuth clients */
        get: operations["listClients"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/internal/clients/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a registered client by ID */
        get: operations["getClient"];
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
        TokenRequest: {
            clientId: string;
            clientSecret: string;
            /** @enum {string} */
            grantType: "client_credentials" | "refresh_token" | "authorization_code";
            scope?: string;
            /** @description Required for refresh_token grant */
            refreshToken?: string;
            /** @description Required for authorization_code grant */
            code?: string;
            /** @description Required for authorization_code grant */
            redirectUri?: string;
        };
        TokenResponse: {
            accessToken: string;
            /** @example Bearer */
            tokenType: string;
            /** @description Seconds until token expires */
            expiresIn: number;
            refreshToken?: string;
            scope?: string;
        };
        IntrospectRequest: {
            token: string;
        };
        IntrospectResponse: {
            active: boolean;
            scope?: string;
            clientId?: string;
            tokenType?: string;
            /** Format: int64 */
            exp?: number;
            /** Format: int64 */
            iat?: number;
            sub?: string;
            aud?: string;
        };
        RevokeRequest: {
            token: string;
            /** @enum {string} */
            token_type_hint?: "access_token" | "refresh_token";
        };
        JWK: {
            /** @example RSA */
            kty: string;
            /** @example sig */
            use: string;
            kid: string;
            /** @example RS256 */
            alg: string;
            /** @description RSA modulus (base64url) */
            n: string;
            /** @description RSA exponent (base64url) */
            e: string;
        };
        JWKS: {
            keys: components["schemas"]["JWK"][];
        };
        /** @description Public view of an OAuth client — secret is never included */
        RegisteredClient: {
            client_id: string;
            name: string;
            allowed_scopes: string[];
        };
        ValidateRequest: {
            token: string;
        };
        ValidateResponse: {
            valid: boolean;
            /** @description Token subject (user ID or service ID) */
            sub?: string;
            scopes?: string[];
            /** @description Present when valid is false */
            error?: string;
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
        /** @description Missing or invalid credentials */
        Unauthorized: {
            headers: {
                [name: string]: unknown;
            };
            content: {
                "application/problem+json": components["schemas"]["Problem"];
            };
        };
        /** @description Resource not found */
        NotFound: {
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
    getJWKS: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description JWKS document */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["JWKS"];
                };
            };
        };
    };
    issueToken: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["TokenRequest"];
            };
        };
        responses: {
            /** @description Token issued */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TokenResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
    authorize: {
        parameters: {
            query: {
                response_type: "code";
                client_id: string;
                redirect_uri: string;
                scope?: string;
                state?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Redirect to client with authorization code */
            302: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            400: components["responses"]["BadRequest"];
        };
    };
    introspectToken: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["IntrospectRequest"];
            };
        };
        responses: {
            /** @description Introspection result — active:false for expired, revoked, or unknown tokens */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["IntrospectResponse"];
                };
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    revokeToken: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["RevokeRequest"];
            };
        };
        responses: {
            /** @description Token revoked (always 200 per RFC 7009 — no error if token is already invalid) */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    validateToken: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ValidateRequest"];
            };
        };
        responses: {
            /** @description Validation result */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ValidateResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
        };
    };
    listClients: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Client list */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["RegisteredClient"][];
                };
            };
        };
    };
    getClient: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Client details */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["RegisteredClient"];
                };
            };
            404: components["responses"]["NotFound"];
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type TokenRequest = components['schemas']['TokenRequest'];
export type TokenResponse = components['schemas']['TokenResponse'];
export type IntrospectRequest = components['schemas']['IntrospectRequest'];
export type IntrospectResponse = components['schemas']['IntrospectResponse'];
export type RevokeRequest = components['schemas']['RevokeRequest'];
export type JWK = components['schemas']['JWK'];
export type JWKS = components['schemas']['JWKS'];
export type RegisteredClient = components['schemas']['RegisteredClient'];
export type ValidateRequest = components['schemas']['ValidateRequest'];
export type ValidateResponse = components['schemas']['ValidateResponse'];
export type Problem = components['schemas']['Problem'];