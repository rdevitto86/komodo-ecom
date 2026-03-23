import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { ShopItemsClient, APIError } from '$lib/server/shop';

export const GET: RequestHandler = async ({ params }) => {
  const client = new ShopItemsClient();
  try {
    const product = await client.getProduct(params.id);
    return json(product);
  } catch (err) {
    if (err instanceof APIError) {
      error(err.status, { message: err.problem.detail ?? err.problem.title });
    }
    error(502, { message: 'upstream error fetching product' });
  }
};
