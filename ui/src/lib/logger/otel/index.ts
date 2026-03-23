import { browser } from '$app/environment';
import { LoggingAdapter } from '../common/adapter';
import type { TelemetryLogEvent } from '../common/schema';
import { getCorrelationId } from '../common/correlation';
import { LOG_ENDPOINT } from '../common/handler';

const SERVICE = import.meta.env.VITE_APP_NAME    || 'komodo-ui';
const VERSION = import.meta.env.VITE_APP_VERSION || 'unknown';
const ENV     = import.meta.env.MODE             || 'development';

/**
 * TelemetryLogger — performance metrics, Core Web Vitals, and span timing.
 *
 * Exported as both TelemetryLogger (preferred) and OtelLogger (legacy alias).
 */
export class TelemetryLogger extends LoggingAdapter {
  constructor() {
    super({
      provider:  'telemetry',
      limit:     10,
      interval:  30_000,    // telemetry is low-urgency, batch aggressively
      endpoint:  LOG_ENDPOINT,
    });
  }

  trace(name: string, attributes?: Partial<TelemetryLogEvent['details']>, requestId?: string): void {
    if (!browser || ENV === 'mock') return;

    const event: TelemetryLogEvent = {
      timestamp:     new Date().toISOString(),
      level:         'info',
      type:          'telemetry',
      service:       SERVICE,
      env:           ENV,
      version:       VERSION,
      requestId,
      correlationId: getCorrelationId(),
      message:       name,
      details: {
        name,
        ...attributes,
      },
    };

    this.send(event);
  }
}

/** @deprecated Use TelemetryLogger */
export const OtelLogger = TelemetryLogger;
