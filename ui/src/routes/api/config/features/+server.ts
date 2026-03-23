import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: features-api Lambda handler not yet implemented.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'Feature flags service is not yet available.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
