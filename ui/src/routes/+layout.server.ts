import type { LayoutServerLoad } from './$types';

// Minimal root layout load — session validation added here when auth-api
// supports the authorization_code flow and user sessions are implemented.
export const load: LayoutServerLoad = async () => {
  return {};
};
