import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { ShopItemsClient, APIError } from '$lib/server/shop';

export const GET: RequestHandler = async () => {
  const client = new ShopItemsClient();
  try {
    const inventory = await client.getInventory();
    return json(inventory);
  } catch (err) {
    if (err instanceof APIError) {
      error(err.status, { message: err.problem.detail ?? err.problem.title });
    }
    error(502, { message: 'upstream error fetching inventory' });
  }
};
