import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: no scheduling service planned yet.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Service scheduling is not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
export const POST: RequestHandler = notImplemented;
