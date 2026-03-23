import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: requires user session (auth-api authorization_code flow not yet implemented).
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'User preferences require an active session. Auth flow not yet deployed.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
export const PUT: RequestHandler = notImplemented;
