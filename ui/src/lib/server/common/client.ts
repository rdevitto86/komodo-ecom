import { isRequestValid } from '$lib/server/utils/http';

type Method = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

export abstract class APIClient {
  #baseURL: string = import.meta.env.API_BASE_URL;

  async send<T>(method: Method, path: string, body?: any, headers?: any): Promise<T> {
    try {
      const url = `${this.#baseURL}${path}`;
      const req: RequestInit = {
        method,
        body: typeof body !== 'string' ? JSON.stringify(body) : body,
        headers: {
          'Content-Type': 'application/json',
          ...headers
        },
      };
      
      if (!isRequestValid(url, req)) {
        const err = new Error('api validation failed');
        logger.error('api request failed validation', err);
        throw err;
      }
      
      const res = await fetch(url, req);

      if (!res.ok) {
        const err = new Error(res.statusText);
        logger.error(`api request failed with status ${res.status}`, err, res);
        throw err;
      }

      return res.json() as Promise<T>;
    } catch (err) {
      logger.error('api request failed', err as Error);
      throw err;
    }
  }
}
