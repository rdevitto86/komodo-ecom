import { logger } from '$lib/logger';

const AUTH_API_URL = process.env.AUTH_API_URL;
if (!AUTH_API_URL) throw new Error('AUTH_API_URL is not defined');

export async function validateApiToken(token: string) {
  try {
    const response = await fetch(`${AUTH_API_URL}/auth/validate`, {
      method: 'POST',
      headers: {
        'Accept': 'application/json;v=1',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (!response.ok) {
      logger.error('Token validation failed', new Error(`Token validation failed: ${response.status}`));
      return null;
    }
    
    const data = await response.json();
    
    return {
      id: data.userId,
      email: data.email,
      isAdmin: data.role === 'admin' || data.isAdmin
    };
  } catch (err) {
    logger.error('Auth API call failed', err as Error);
    return null;
  }
}