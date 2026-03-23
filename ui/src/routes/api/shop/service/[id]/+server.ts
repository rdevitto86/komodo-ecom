import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: service catalog endpoint not yet implemented in shop-items-api.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Service catalog is not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
