import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { ShopItemsClient, APIError } from '$lib/server/shop';

export const load: PageServerLoad = async ({ params }) => {
  const client = new ShopItemsClient();
  try {
    const product = await client.getProduct(params.id);
    return { product };
  } catch (err) {
    if (err instanceof APIError) {
      error(err.status, { message: err.problem.detail ?? err.problem.title });
    }
    error(502, { message: 'failed to load product' });
  }
};
