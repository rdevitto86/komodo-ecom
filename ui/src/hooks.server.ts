import { type Handle } from '@sveltejs/kit';
import { serverLogger } from '$lib/server/logging';
import { randomUUID } from 'crypto';

export const handle: Handle = async ({ event, resolve }) => {
  const requestId = event.request.headers.get('x-request-id') ?? randomUUID();
  event.locals.requestId = requestId;

  const start = Date.now();
  const { pathname } = event.url;

  const response = await resolve(event, {
    transformPageChunk: ({ html }) => html,
  });

  // Skip logging static assets to avoid noise
  const isAsset = /\.(jpg|jpeg|png|svg|webp|gif|ico|css|js|woff2?)$/.test(pathname);
  if (!isAsset) {
    serverLogger.info('request', {
      requestId,
      method: event.request.method,
      path: pathname,
      status: response.status,
      durationMs: Date.now() - start,
    });
  }

  response.headers.set('x-request-id', requestId);
  return response;
};
