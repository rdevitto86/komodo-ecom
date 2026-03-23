import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: no shipping service planned yet.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Shipping estimates are not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
