import { json, error, isHttpError } from '@sveltejs/kit';
import { getPageContentFromS3 } from '$lib/server/s3';
import { logger } from '$lib/logger';

export async function GET({ params }) {
  try {
    const product = await getPageContentFromS3(`product-${params.id}`);
    if (!product) throw error(404, `Product ${params.id} not found`);
    return json(product);
  } catch (err) {
    if (isHttpError(err)) {
      logger.error(err?.body?.message);
      throw err;
    }
    logger.error('Failed to fetch product');
    throw error(500, 'Failed to fetch product');
  }
}