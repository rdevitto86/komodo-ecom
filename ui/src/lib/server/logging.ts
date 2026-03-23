import type { BaseLogEvent, LogLevel } from '$lib/logger/common/schema';

const SERVICE = process.env.APP_NAME    || 'komodo-ui';
const ENV     = process.env.NODE_ENV    || 'development';
const VERSION = process.env.APP_VERSION || 'unknown';

const LEVEL_WEIGHT: Record<LogLevel, number> = {
  debug: 1, info: 2, warn: 3, error: 4
};

const minLevel: number = LEVEL_WEIGHT[
  ((process.env.LOG_LEVEL || (ENV === 'production' ? 'warn' : 'info')) as LogLevel)
] ?? LEVEL_WEIGHT.info;

const isLocal = ENV === 'development' || ENV === 'local';

// --- Console string format (matches Go SDK KomodoTextHandler) ---
// Local:  2026-03-22T14:23:01Z [ERROR] req-abc123 | message | key=val
// Prod:   structured JSON to stdout (CloudWatch picks it up)

function toLogfmt(obj: Record<string, unknown>): string {
  return Object.entries(obj)
    .filter(([, v]) => v !== undefined && v !== null)
    .map(([k, v]) => {
      const s = typeof v === 'string' ? v : JSON.stringify(v);
      return /[\s"=]/.test(s) ? `${k}="${s.replace(/"/g, '\\"')}"` : `${k}=${s}`;
    })
    .join(' ');
}

function formatString(event: BaseLogEvent): string {
  const lvl    = `[${event.level.toUpperCase()}]`;
  const reqId  = event.requestId || '-';
  const detail = event.details ? ' | ' + toLogfmt(event.details as Record<string, unknown>) : '';
  return `${event.timestamp} ${lvl} ${reqId} | ${event.message}${detail}`;
}

/**
 * logEvent — server-side logger for proxied frontend events.
 *
 * In production (adapter-node on EC2/Fargate): prints structured JSON to
 * stdout. The CloudWatch agent on the host ships it to the configured log group.
 *
 * Locally: prints the same human-readable string format as the Go SDK.
 */
export function logEvent(event: BaseLogEvent): void {
  const weight = LEVEL_WEIGHT[event.level] ?? 0;
  if (weight < minLevel) return;

  // Enrich with server-known fields if the client didn't supply them
  const enriched: BaseLogEvent = {
    ...event,
    service: event.service || SERVICE,
    env:     event.env     || ENV,
    version: event.version || VERSION,
  };

  if (isLocal) {
    const fn = enriched.level === 'error' ? console.error
      : enriched.level === 'warn'  ? console.warn
      : console.log;
    fn(formatString(enriched));
  } else {
    // Structured JSON — one line per event for CW Insights
    console.log(JSON.stringify(enriched));
  }
}

/**
 * logEvents — batch variant used by the /api/log route.
 */
export function logEvents(events: BaseLogEvent[]): void {
  for (const e of events) {
    logEvent(e);
  }
}

/**
 * serverLogger — structured server-side logger for SvelteKit load functions,
 * hooks, and +server.ts routes. Matches the same field set as the Go SDK.
 */
export const serverLogger = {
  debug: (message: string, details?: Record<string, unknown>, requestId?: string) =>
    logEvent(makeEvent('debug', message, details, requestId)),
  info: (message: string, details?: Record<string, unknown>, requestId?: string) =>
    logEvent(makeEvent('info', message, details, requestId)),
  warn: (message: string, details?: Record<string, unknown>, requestId?: string) =>
    logEvent(makeEvent('warn', message, details, requestId)),
  error: (message: string, details?: Record<string, unknown>, requestId?: string) =>
    logEvent(makeEvent('error', message, details, requestId)),
};

function makeEvent(
  level: LogLevel,
  message: string,
  details?: Record<string, unknown>,
  requestId?: string
): BaseLogEvent {
  return {
    timestamp:     new Date().toISOString(),
    level,
    type:          'runtime',
    service:       SERVICE,
    env:           ENV,
    version:       VERSION,
    requestId,
    correlationId: 'server',
    message,
    details,
  };
}
