// Catalog types — mirrors komodo-forge-sdk-ts/src/shared/types/products.ts
// and the Go pkg/v1/models shapes in komodo-shop-items-api.
// Keep in sync with both when models evolve.

export type StockCode = 'IS' | 'OS' | 'LS' | 'PO' | 'SO' | 'BO' | 'DC' | 'TU';

export interface Product {
  id: string;
  slug: string;
  name: string;
  description: string;
  brand?: string;
  manufacturer?: string;
  status: 'active' | 'draft' | 'archived';
  currency?: string;
  price?: number;
  compareAtPrice?: number;
  cost?: number;
  taxCode?: string;
  trackInventory: boolean;
  minOrderQuantity?: number;
  maxOrderQuantity?: number;
  customizationOptions?: CustomizationOption[];
  addOns?: AddOn[];
  relatedProductIds?: string[];
  variants: Variant[];
  specs?: Record<string, unknown>;
  meta?: ProductMeta;
  seo?: SEO;
  createdAt?: string;
  updatedAt?: string;
}

export interface ProductMeta {
  tags?: string[];
  categories?: string[];
  isPopular?: boolean;
  isFeatured?: boolean;
  isNew?: boolean;
  [key: string]: unknown;
}

export interface SEO {
  title?: string;
  description?: string;
  keywords?: string[];
}

export interface Variant {
  id: string;
  sku?: string;
  upc?: string;
  gtin?: string;
  ean?: string;
  model?: string;
  name: string;
  description?: string;
  price: number;
  compareAtPrice?: number;
  cost?: number;
  taxCode?: string;
  stockQty?: number;
  stockCode?: StockCode;
  images?: ProductImage[];
  optionCombination?: Record<string, string>;
  weight?: number;
  weightUnit?: 'lb' | 'kg' | 'oz' | 'g';
  dimensions?: { length?: number; width?: number; height?: number; unit?: 'in' | 'cm' };
  requiresShipping?: boolean;
  shippingClass?: string;
  isDefault?: boolean;
}

export interface ProductImage {
  url: string;
  alt?: string;
  isPrimary?: boolean;
  variantIds?: string[];
  type?: 'image' | 'video' | 'spin360' | 'model3d';
  spin360?: { frames: string[]; frameCount: number; startFrame?: number };
  model3d?: { modelUrl: string; format: string; textureUrls?: string[]; thumbnailUrl?: string };
}

export interface CustomizationOption {
  id: string;
  name: string;
  type: string;
  required: boolean;
  displayOrder: number;
  values: CustomizationValue[];
}

export interface CustomizationValue {
  id: string;
  label: string;
  value: string;
  priceModifier?: number;
  hexColor?: string;
  imageUrl?: string;
  stockCode?: StockCode;
  stockQty?: number;
  isDefault?: boolean;
  disabled?: boolean;
  disabledReason?: string;
}

export interface AddOn {
  id: string;
  sku?: string;
  name: string;
  description?: string;
  manufacturer?: string;
  price: number;
  compareAtPrice?: number;
  imageUrl?: string;
  stockCode?: StockCode;
  stockQty?: number;
  weight?: number;
  requiresShipping?: boolean;
  maxQuantity?: number;
  isRecommended?: boolean;
  compatibleWith?: { optionIds?: string[]; variantIds?: string[] };
}

export interface InventoryItem {
  sku: string;
  name: string;
  stockQty: number;
  stockCode: StockCode;
  price: number;
  updatedAt?: string;
}

export interface InventoryResponse {
  items: InventoryItem[];
  total: number;
}
