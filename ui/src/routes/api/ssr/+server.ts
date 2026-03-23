import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// Blocked: ssr-engine-svelte not yet active.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'SSR fragment rendering is not yet active.' },
    { status: 501 },
  );

export const POST: RequestHandler = notImplemented;
