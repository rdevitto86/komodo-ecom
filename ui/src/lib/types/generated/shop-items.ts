// Auto-generated from apis/komodo-shop-items-api/docs/openapi.yaml
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
    "/item/inventory": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get inventory summary for all SKUs */
        get: operations["getInventory"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/item/{sku}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get a product or service by SKU */
        get: operations["getItemBySKU"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/item/suggestion": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Get personalized product suggestions */
        post: operations["getSuggestions"];
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
        /**
         * @description Inventory availability code:
         *     - IS: In Stock
         *     - OS: Out of Stock
         *     - LS: Limited Supply
         *     - PO: Pre-Order
         *     - SO: Sold Out
         *     - BO: Backorder
         *     - DC: Discontinued
         *     - TU: Temporarily Unavailable
         * @enum {string}
         */
        StockCode: "IS" | "OS" | "LS" | "PO" | "SO" | "BO" | "DC" | "TU";
        Dimensions: {
            /** Format: double */
            length?: number;
            /** Format: double */
            width?: number;
            /** Format: double */
            height?: number;
            /** @enum {string} */
            unit?: "in" | "cm";
        };
        Spin360: {
            frames: string[];
            frameCount: number;
            startFrame?: number;
        };
        Model3D: {
            /** Format: uri */
            modelUrl: string;
            format: string;
            textureUrls?: string[];
            /** Format: uri */
            thumbnailUrl?: string;
        };
        ProductImage: {
            /** Format: uri */
            url: string;
            alt?: string;
            isPrimary?: boolean;
            variantIds?: string[];
            /** @enum {string} */
            type?: "image" | "video" | "spin360" | "model3d";
            spin360?: components["schemas"]["Spin360"];
            model3d?: components["schemas"]["Model3D"];
        };
        CustomizationValue: {
            id: string;
            label: string;
            value: string;
            /** Format: double */
            priceModifier?: number;
            hexColor?: string;
            /** Format: uri */
            imageUrl?: string;
            stockCode?: components["schemas"]["StockCode"];
            stockQty?: number;
            isDefault?: boolean;
            disabled?: boolean;
            disabledReason?: string;
        };
        CustomizationOption: {
            id: string;
            name: string;
            type: string;
            required: boolean;
            displayOrder: number;
            values: components["schemas"]["CustomizationValue"][];
        };
        Compatibility: {
            optionIds?: string[];
            variantIds?: string[];
        };
        AddOn: {
            id: string;
            sku?: string;
            name: string;
            description?: string;
            manufacturer?: string;
            /** Format: double */
            price: number;
            /** Format: double */
            compareAtPrice?: number;
            /** Format: uri */
            imageUrl?: string;
            stockCode?: components["schemas"]["StockCode"];
            stockQty?: number;
            /** Format: double */
            weight?: number;
            requiresShipping?: boolean;
            maxQuantity?: number;
            isRecommended?: boolean;
            compatibleWith?: components["schemas"]["Compatibility"];
        };
        Variant: {
            id: string;
            sku?: string;
            upc?: string;
            gtin?: string;
            ean?: string;
            model?: string;
            name: string;
            description?: string;
            /** Format: double */
            price: number;
            /** Format: double */
            compareAtPrice?: number;
            /** Format: double */
            cost?: number;
            taxCode?: string;
            stockQty?: number;
            stockCode?: components["schemas"]["StockCode"];
            images?: components["schemas"]["ProductImage"][];
            optionCombination?: {
                [key: string]: string;
            };
            /** Format: double */
            weight?: number;
            /** @enum {string} */
            weightUnit?: "lb" | "kg" | "oz" | "g";
            dimensions?: components["schemas"]["Dimensions"];
            requiresShipping?: boolean;
            shippingClass?: string;
            isDefault?: boolean;
        };
        ProductMeta: {
            tags?: string[];
            categories?: string[];
            isPopular?: boolean;
            isFeatured?: boolean;
            isNew?: boolean;
            extra?: {
                [key: string]: unknown;
            };
        };
        SEO: {
            title?: string;
            description?: string;
            keywords?: string[];
        };
        Product: {
            id: string;
            slug: string;
            name: string;
            description: string;
            brand?: string;
            manufacturer?: string;
            /** @enum {string} */
            status: "active" | "draft" | "archived";
            currency?: string;
            /** Format: double */
            price?: number;
            /** Format: double */
            compareAtPrice?: number;
            /** Format: double */
            cost?: number;
            taxCode?: string;
            trackInventory: boolean;
            minOrderQuantity?: number;
            maxOrderQuantity?: number;
            customizationOptions?: components["schemas"]["CustomizationOption"][];
            addOns?: components["schemas"]["AddOn"][];
            relatedProductIds?: string[];
            variants: components["schemas"]["Variant"][];
            specs?: {
                [key: string]: unknown;
            };
            meta?: components["schemas"]["ProductMeta"];
            seo?: components["schemas"]["SEO"];
            /** Format: date-time */
            createdAt?: string;
            /** Format: date-time */
            updatedAt?: string;
        };
        TimeSlot: {
            /** @description HH:MM (24h) */
            start: string;
            /** @description HH:MM (24h) */
            end: string;
        };
        ServiceAreaRadius: {
            /** Format: double */
            miles: number;
            centerZip: string;
        };
        ServiceAvailability: {
            daysOfWeek?: number[];
            timeSlots?: components["schemas"]["TimeSlot"][];
            blackoutDates?: string[];
            leadTimeDays?: number;
            maxBookingsPerDay?: number;
            serviceAreaZipCodes?: string[];
            serviceAreaRadius?: components["schemas"]["ServiceAreaRadius"];
        };
        ServiceDuration: {
            estimated: number;
            /** @enum {string} */
            unit: "minutes" | "hours" | "days";
            min?: number;
            max?: number;
        };
        ServiceImage: {
            /** Format: uri */
            url: string;
            alt?: string;
            isPrimary?: boolean;
            type?: string;
        };
        ServiceOption: {
            id: string;
            name: string;
            description?: string;
            type: string;
            /** Format: double */
            priceModifier: number;
            /** Format: double */
            durationModifier?: number;
            maxBookings?: number;
            isDefault?: boolean;
            requiresProducts?: string[];
            compatibleWith?: string[];
        };
        Service: {
            id: string;
            slug: string;
            sku?: string;
            name: string;
            description: string;
            category: string;
            /** @enum {string} */
            status: "active" | "draft" | "archived";
            /** Format: double */
            price: number;
            /** Format: double */
            compareAtPrice?: number;
            /** Format: double */
            cost?: number;
            currency?: string;
            taxCode?: string;
            duration?: components["schemas"]["ServiceDuration"];
            images?: components["schemas"]["ServiceImage"][];
            serviceOptions?: components["schemas"]["ServiceOption"][];
            locationTypes: string[];
            availability?: components["schemas"]["ServiceAvailability"];
            requirements?: string[];
            includedItems?: string[];
            relatedServiceIds?: string[];
            relatedProductIds?: string[];
            meta?: components["schemas"]["ProductMeta"];
            seo?: components["schemas"]["SEO"];
            /** Format: date-time */
            createdAt?: string;
            /** Format: date-time */
            updatedAt?: string;
        };
        InventoryItem: {
            sku: string;
            name: string;
            stockQty: number;
            stockCode: components["schemas"]["StockCode"];
            /** Format: double */
            price: number;
            /** Format: date-time */
            updatedAt?: string;
        };
        InventoryResponse: {
            items: components["schemas"]["InventoryItem"][];
            total: number;
        };
        SuggestionRequest: {
            userId?: string;
            /** @description SKUs of items currently in cart or recently viewed */
            skus?: string[];
            categories?: string[];
            /** @default 10 */
            limit: number;
        };
        SuggestionResponse: {
            suggestions: components["schemas"]["Product"][];
            total: number;
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
        /** @description Item not found */
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
    getInventory: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Inventory list */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["InventoryResponse"];
                };
            };
        };
    };
    getItemBySKU: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                sku: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Product or service details */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Product"] | components["schemas"]["Service"];
                };
            };
            404: components["responses"]["NotFound"];
        };
    };
    getSuggestions: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["SuggestionRequest"];
            };
        };
        responses: {
            /** @description Suggested products */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SuggestionResponse"];
                };
            };
            400: components["responses"]["BadRequest"];
            401: components["responses"]["Unauthorized"];
        };
    };
}


// Named schema exports — import directly without referencing components['schemas']
export type StockCode = components['schemas']['StockCode'];
export type Dimensions = components['schemas']['Dimensions'];
export type Spin360 = components['schemas']['Spin360'];
export type Model3D = components['schemas']['Model3D'];
export type ProductImage = components['schemas']['ProductImage'];
export type CustomizationValue = components['schemas']['CustomizationValue'];
export type CustomizationOption = components['schemas']['CustomizationOption'];
export type Compatibility = components['schemas']['Compatibility'];
export type AddOn = components['schemas']['AddOn'];
export type Variant = components['schemas']['Variant'];
export type ProductMeta = components['schemas']['ProductMeta'];
export type SEO = components['schemas']['SEO'];
export type Product = components['schemas']['Product'];
export type TimeSlot = components['schemas']['TimeSlot'];
export type ServiceAreaRadius = components['schemas']['ServiceAreaRadius'];
export type ServiceAvailability = components['schemas']['ServiceAvailability'];
export type ServiceDuration = components['schemas']['ServiceDuration'];
export type ServiceImage = components['schemas']['ServiceImage'];
export type ServiceOption = components['schemas']['ServiceOption'];
export type Service = components['schemas']['Service'];
export type InventoryItem = components['schemas']['InventoryItem'];
export type InventoryResponse = components['schemas']['InventoryResponse'];
export type SuggestionRequest = components['schemas']['SuggestionRequest'];
export type SuggestionResponse = components['schemas']['SuggestionResponse'];
export type Problem = components['schemas']['Problem'];