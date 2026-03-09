import { error, type Handle } from '@sveltejs/kit';
import { validateApiToken } from '$lib/server/auth';
import { logger } from '$lib/logger';
import { initCacheWriter, getFrequentRoutes, getCacheStats } from '$lib/server/cache';
import { fetchAndCacheFromS3 } from '$lib/server/s3';

// ── Startup: init writer thread → warm caches ──────────────────────────────
(async () => {
  logger.info('[Warmup] Starting cache writer thread...');
  await initCacheWriter();

  const PRELOADED_TEMPLATES = ['products', 'services'];

  // 1. Pre-load common templates from S3
  const preloadResults = await Promise.allSettled(
    PRELOADED_TEMPLATES.map(key => fetchAndCacheFromS3(key, true))
  );
  const succeeded = preloadResults.filter(result => result.status === 'fulfilled').length;
  logger.info(`[Warmup] Preloaded ${succeeded}/${PRELOADED_TEMPLATES.length} common templates`);

  // 2. Re-warm previously promoted frequent routes
  const frequentRoutes = await getFrequentRoutes();
  if (frequentRoutes.length > 0) {
    const freqResults = await Promise.allSettled(
      frequentRoutes.map(key => fetchAndCacheFromS3(key, false))
    );
    const freqSucceeded = freqResults.filter(result => result.status === 'fulfilled').length;
    logger.info(`[Warmup] Re-warmed ${freqSucceeded}/${frequentRoutes.length} frequent templates`);
  }

  const stats = await getCacheStats();
  logger.info('[Warmup] Cache stats:', stats);
})().catch(err => logger.error('[Warmup] Failed:', err as Error));

const PROTECTED_ROUTES = [
  '/orders',
  '/marketing/*',
  '/services/scheduling'
].map(route => new RegExp(`/api(/v\\d+)?${route.replace(/\*/g, '(/.*)?')}(/|$)`));

/**
 * TODO: Document this
 * @param param0 
 * @returns 
 */
export const handle: Handle = async ({ event, resolve }) => {
  try {
    const start = Date.now();
    const path = event.url.pathname;
    const token = event.request.headers.get('Authorization')?.replace('Bearer ', '');
    
    if (token) {
      event.locals.user = await validateApiToken(token);
    }
    
    if (PROTECTED_ROUTES.some(pattern => pattern.test(path)) && !event.locals.user) {
      throw error(401, 'Unauthorized - Valid token required');
    }
    
    const response = await resolve(event);
    
    logger.info('Request completed', {
      method: event.request.method,
      path,
      status: response.status,
      duration: Date.now() - start,
      userId: event.locals.user?.id
    });
    
    return response;
  } catch (err) {
    logger.error('Request failed', err as Error);
    throw err;
  }
};
