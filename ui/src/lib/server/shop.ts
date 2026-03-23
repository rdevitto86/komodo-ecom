import { APIClient, APIError } from './common/client';
import { getServerConfig } from './secrets';
import type { Product, InventoryResponse } from '$lib/types/catalog';

export class ShopItemsClient extends APIClient {
  constructor() {
    super(getServerConfig().shopItemsApiUrl);
  }

  async getProduct(sku: string): Promise<Product> {
    return this.send<Product>('GET', `/item/${encodeURIComponent(sku)}`);
  }

  async getInventory(): Promise<InventoryResponse> {
    return this.send<InventoryResponse>('GET', '/item/inventory');
  }
}

export { APIError };
