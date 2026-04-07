// Auto-generated from apis/komodo-order-returns-api/docs/openapi.yaml
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
    "/me/returns": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List user's return requests */
        get: operations["listMyReturns"];
        put?: never;
        /** Initiate a return request */
        post: operations["createReturn"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/me/returns/{returnId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get return status */
        get: operations["getMyReturn"];
        put?: never;
        post?: never;
        /** Cancel a pending return request */
        delete: operations["cancelMyReturn"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/returns/{returnId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get return (admin/service use) */
        get: operations["getReturn"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/returns/{returnId}/approve": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        /** Approve return and trigger refund */
        put: operations["approveReturn"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/returns/{returnId}/receive": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        /** Mark items received, trigger restock and loyalty reversal */
        put: operations["receiveReturn"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/returns/{returnId}/reject": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        /** Reject return request */
        put: operations["rejectReturn"];
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
        /** @enum {string} */
        ReturnStatus: "requested" | "approved" | "received" | "processed" | "rejected" | "cancelled";
        ReturnItem: {
            item_id: string;
            sku: string;
            quantity: number;
            /** @enum {string} */
            reason: "defective" | "wrong_item" | "not_as_described" | "changed_mind" | "other";
            notes?: string;
        };
        Return: {
            /** Format: uuid */
            return_id: string;
            order_id: string;
            status: components["schemas"]["ReturnStatus"];
            items: components["schemas"]["ReturnItem"][];
            refund_amount_cents?: number;
            /** @description Set when refund is initiated */
            refund_id?: string;
            /** Format: date-time */
            return_window_expires?: string;
            rejection_reason?: string;
            /** Format: date-time */
            created_at: string;
            /** Format: date-time */
            updated_at: string;
        };
        ReturnList: {
            returns: components["schemas"]["Return"][];
            next_cursor?: string;
        };
        CreateReturnRequest: {
            order_id: string;
            items: components["schemas"]["ReturnItem"][];
        };
        ApproveReturnRequest: {
            /** @description Override refund amount. Defaults to full item value if omitted. */
            refund_amount_cents?: number;
        };
        RejectReturnRequest: {
            reason: string;
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
        /** @description Return not found */
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
        ReturnId: string;
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
    listMyReturns: {
        parameters: {
            query?: {
                status?: components["schemas"]["ReturnStatus"];
                limit?: number;
                cursor?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Return requests */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ReturnList"];
                };
            };
            401: components["responses"]["Unauthorized"];
        };
    };
    createReturn: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateReturnRequest"];
            };
        };
        responses: {
            /** @description Return request created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Return"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
            /** @description Order not eligible for return (outside window, already returned, or not fulfilled) */
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
    getMyReturn: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                returnId: components["parameters"]["ReturnId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Return details */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Return"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    cancelMyReturn: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                returnId: components["parameters"]["ReturnId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Return cancelled */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
            /** @description Return cannot be cancelled (already approved or processed) */
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
    getReturn: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                returnId: components["parameters"]["ReturnId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Return details */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Return"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
    approveReturn: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                returnId: components["parameters"]["ReturnId"];
            };
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": components["schemas"]["ApproveReturnRequest"];
            };
        };
        responses: {
            /** @description Return approved */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Return"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
            /** @description Return not in approvable state */
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
    receiveReturn: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                returnId: components["parameters"]["ReturnId"];
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Return processed */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Return"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
            /** @description Return not in approvable state */
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
    rejectReturn: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                returnId: components["parameters"]["ReturnId"];
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["RejectReturnRequest"];
            };
        };
        responses: {
            /** @description Return rejected */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Return"];
                };
            };
            401: components["responses"]["Unauthorized"];
            404: components["responses"]["NotFound"];
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type ReturnStatus = components['schemas']['ReturnStatus'];
export type ReturnItem = components['schemas']['ReturnItem'];
export type Return = components['schemas']['Return'];
export type ReturnList = components['schemas']['ReturnList'];
export type CreateReturnRequest = components['schemas']['CreateReturnRequest'];
export type ApproveReturnRequest = components['schemas']['ApproveReturnRequest'];
export type RejectReturnRequest = components['schemas']['RejectReturnRequest'];
export type Problem = components['schemas']['Problem'];