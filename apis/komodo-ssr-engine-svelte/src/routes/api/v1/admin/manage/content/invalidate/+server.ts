import { json, error } from '@sveltejs/kit';
import { invalidatePage } from '$lib/server/cloudfront';
import type { RequestHandler } from '../../$types';

export const POST: RequestHandler = async ({ request, locals }) => {
  // TODO: Add authentication check
  // if (!locals.user?.isAdmin) {
  //   throw error(403, 'Forbidden');
  // }

  const { pageKey, cloudfront = true } = await request.json();

  if (!pageKey) {
    throw error(400, 'pageKey required');
  }

  if (cloudfront) {
    try {
      await invalidatePage(pageKey);
    } catch (err) {
      console.error('CloudFront invalidation failed:', err);
    }
  }

  return json({
    success: true,
    invalidated: pageKey,
    timestamp: new Date().toISOString()
  });
};
