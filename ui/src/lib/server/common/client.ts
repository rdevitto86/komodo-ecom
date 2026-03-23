import { serverLogger } from '$lib/server/logging';
import { isRequestValid } from '$lib/server/utils/http';

type Method = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

// RFC 7807 Problem+JSON — shape returned by all Go services on error.
export interface ProblemDetail {
  type?: string;
  title: string;
  status: number;
  detail?: string;
  instance?: string;
}

export class APIError extends Error {
  constructor(
    public readonly status: number,
    public readonly problem: ProblemDetail,
  ) {
    super(problem.detail ?? problem.title);
    this.name = 'APIError';
  }
}

export abstract class APIClient {
  constructor(protected readonly baseURL: string) {}

  async send<T>(method: Method, path: string, body?: unknown, headers?: Record<string, string>): Promise<T> {
    const url = `${this.baseURL}${path}`;
    const req: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...headers,
      },
      ...(body !== undefined && { body: JSON.stringify(body) }),
    };

    if (!isRequestValid(url, req)) {
      const err = new Error(`invalid request: ${method} ${url}`);
      serverLogger.error('api request failed validation', { method, path });
      throw err;
    }

    const res = await fetch(url, req);

    if (!res.ok) {
      let problem: ProblemDetail;
      try {
        problem = await res.json() as ProblemDetail;
      } catch {
        problem = { title: res.statusText, status: res.status };
      }
      serverLogger.error('api request failed', { method, path, status: res.status, detail: problem.detail });
      throw new APIError(res.status, problem);
    }

    return res.json() as Promise<T>;
  }
}
