import { env } from '$env/dynamic/private';

/**
 * Runtime config resolved from environment variables.
 * On EC2/Fargate these are injected by the container runtime.
 * Locally they come from .env files or docker-compose environment blocks.
 *
 * NOTE: Never expose secrets through a BFF API route.
 * All secret access stays server-side within this module.
 */
export function getServerConfig() {
  return {
    env: env.ENV ?? 'development',
    authApiUrl:      env.AUTH_API_URL      ?? 'http://localhost:7011',
    userApiUrl:      env.USER_API_URL      ?? 'http://localhost:7051',
    shopItemsApiUrl: env.SHOP_ITEMS_API_URL ?? 'http://localhost:7041',
    cartApiUrl:      env.CART_API_URL      ?? 'http://localhost:7043',
  };
}
