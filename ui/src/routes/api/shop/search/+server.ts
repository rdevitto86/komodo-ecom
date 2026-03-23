import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: search-api not yet deployed.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Search is not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
