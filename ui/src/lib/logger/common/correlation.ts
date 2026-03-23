import { browser } from '$app/environment';

const CORRELATION_KEY = 'komodo_cid';

let _cached: string | null = null;

/**
 * Returns the browser's correlation ID — a short-lived fingerprint that
 * identifies a single browser session (tab lifetime, not across tabs).
 *
 * Generated once on first call via crypto.randomUUID(), stored in sessionStorage.
 * Survives page navigations within the tab; resets when the tab closes.
 *
 * SSR context: returns 'server' — correlation IDs are a client-only concept.
 */
export function getCorrelationId(): string {
  if (!browser) return 'server';
  if (_cached) return _cached;

  try {
    const stored = sessionStorage.getItem(CORRELATION_KEY);
    if (stored) {
      _cached = stored;
      return _cached;
    }
    const id = crypto.randomUUID();
    sessionStorage.setItem(CORRELATION_KEY, id);
    _cached = id;
    return _cached;
  } catch {
    // sessionStorage unavailable (private browsing restrictions, etc.)
    _cached = crypto.randomUUID();
    return _cached;
  }
}
