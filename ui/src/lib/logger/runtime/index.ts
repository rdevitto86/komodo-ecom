import { browser } from '$app/environment';
import type { RuntimeLogEvent } from '../common/schema';
import { getCorrelationId } from '../common/correlation';
import { redact } from '../common/redaction';
import { WorkerHandler, LOG_ENDPOINT } from '../common/handler';

export type LogLevel = 'off' | 'debug' | 'info' | 'warn' | 'error';

const LEVEL_WEIGHT: Record<LogLevel, number> = {
  off: 0, debug: 1, info: 2, warn: 3, error: 4
};

const SERVICE  = import.meta.env.VITE_APP_NAME    || 'komodo-ui';
const VERSION  = import.meta.env.VITE_APP_VERSION || 'unknown';
const ENV      = import.meta.env.MODE             || 'development';

// Verbose console format in non-production envs
const verbose = import.meta.env.DEV
  || ENV === 'development'
  || ENV === 'qa'
  || ENV === 'staging';

// Remote log shipping requires a live server (adapter-node, not static/mock)
const remoteEnabled = browser && ENV !== 'mock';

let activeLevel: number = LEVEL_WEIGHT[
  ((import.meta.env.VITE_LOG_LEVEL || (import.meta.env.PROD ? 'error' : 'warn')) as LogLevel)
] ?? LEVEL_WEIGHT.warn;

// Configure the worker once for the runtime provider
if (remoteEnabled) {
  WorkerHandler.init();
  WorkerHandler.config('runtime', { limit: 10, interval: 10_000, endpoint: LOG_ENDPOINT });
}

// --- Console format ---
// PROD:     [ERROR] message
// DEV/QA:   2026-03-22T14:23:01Z [ERROR] req-abc123 | message | key=val key2="val val"

function toLogfmt(obj: Record<string, unknown>): string {
  return Object.entries(obj)
    .filter(([, v]) => v !== undefined && v !== null)
    .map(([k, v]) => {
      const s = typeof v === 'string' ? v : JSON.stringify(v);
      return /[\s"=]/.test(s) ? `${k}="${s.replace(/"/g, '\\"')}"` : `${k}=${s}`;
    })
    .join(' ');
}

function consoleFormat(event: RuntimeLogEvent): string {
  const lvl = `[${event.level.toUpperCase()}]`;
  if (!verbose) return `${lvl} ${event.message}`;

  const reqId  = event.requestId || '-';
  const detail = event.details ? ' | ' + toLogfmt(event.details as Record<string, unknown>) : '';
  return `${event.timestamp} ${lvl} ${reqId} | ${event.message}${detail}`;
}

function buildEvent(
  level: Exclude<LogLevel, 'off'>,
  message: string,
  details?: RuntimeLogEvent['details'],
  requestId?: string
): RuntimeLogEvent {
  return {
    timestamp:     new Date().toISOString(),
    level,
    type:          'runtime',
    service:       SERVICE,
    env:           ENV,
    version:       VERSION,
    requestId,
    correlationId: getCorrelationId(),
    message,
    details,
  };
}

const logger = {
  init(config: { level?: LogLevel }) {
    if (config.level) activeLevel = LEVEL_WEIGHT[config.level] ?? activeLevel;
  },

  debug(message: string, details?: RuntimeLogEvent['details'], requestId?: string) {
    if (activeLevel > LEVEL_WEIGHT.debug) return;
    const event = buildEvent('debug', message, details, requestId);
    if (verbose) console.debug(consoleFormat(event));
  },

  info(message: string, details?: RuntimeLogEvent['details'], requestId?: string) {
    if (activeLevel > LEVEL_WEIGHT.info) return;
    const event = buildEvent('info', message, details, requestId);
    if (verbose) console.info(consoleFormat(event));
    if (remoteEnabled) WorkerHandler.send('runtime', redact(event));
  },

  warn(message: string, details?: RuntimeLogEvent['details'], requestId?: string) {
    if (activeLevel > LEVEL_WEIGHT.warn) return;
    const event = buildEvent('warn', message, details, requestId);
    console.warn(consoleFormat(event));
    if (remoteEnabled) WorkerHandler.send('runtime', redact(event));
  },

  error(message: string, details?: RuntimeLogEvent['details'], requestId?: string) {
    if (activeLevel > LEVEL_WEIGHT.error) return;
    const event = buildEvent('error', message, details, requestId);
    console.error(consoleFormat(event));
    if (remoteEnabled) WorkerHandler.send('runtime', redact(event));
  },
};

export default logger;
