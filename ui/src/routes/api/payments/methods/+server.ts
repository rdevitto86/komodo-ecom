import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: payments-api Lambda handler not yet implemented.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Payment methods are not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
export const POST: RequestHandler = notImplemented;
export const DELETE: RequestHandler = notImplemented;
