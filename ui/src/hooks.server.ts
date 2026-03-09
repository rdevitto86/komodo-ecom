import { error, type Handle } from '@sveltejs/kit';
import { SecretsAPI } from '$lib/server/secrets';
import { Logger } from '$lib/logger/splunk';

// Bootrap globals and secrets
(() => async () => {
  try {
    const secretsAPI = new SecretsAPI();
    await secretsAPI.getSecrets();
  } catch (err) {
    console.error('❌ Failed to load secrets:', err);
  }

  globalThis.logger = new Logger();
})();

export const handle: Handle = async ({ event, resolve }) => {
  try {
    const isStaticAsset = event.url.pathname.match(/\.(jpg|jpeg|png|svg|webp|gif|css|js|woff2|ico)$/);
    if (isStaticAsset) return await resolve(event);
  } catch (err) {
    logger.error('❌ Failed to load secrets:', err as Error);
    error(500, 'Internal Server Configuration Error');
  }
  return await resolve(event);
};
