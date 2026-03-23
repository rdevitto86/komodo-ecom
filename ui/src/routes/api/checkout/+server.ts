import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: requires cart-api implementation and payments-api Lambda handler.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Checkout is not yet available.' },
    { status: 501 },
  );

export const POST: RequestHandler = notImplemented;
