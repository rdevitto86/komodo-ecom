import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: support-api not yet implemented.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Live chat support is not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
export const POST: RequestHandler = notImplemented;
