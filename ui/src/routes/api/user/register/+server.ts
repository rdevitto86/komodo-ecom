import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: requires auth-api authorization_code flow before user-api profile creation.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'User registration is not yet available.' },
    { status: 501 },
  );

export const POST: RequestHandler = notImplemented;
