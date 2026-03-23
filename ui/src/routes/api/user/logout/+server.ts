import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: auth-api authorization_code grant not yet implemented.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'User logout is not yet available. Requires auth-api token revocation flow.' },
    { status: 501 },
  );

export const POST: RequestHandler = notImplemented;
