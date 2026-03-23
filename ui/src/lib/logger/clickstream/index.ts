import { browser } from '$app/environment';
import { LoggingAdapter } from '../common/adapter';
import type { ClickstreamLogEvent } from '../common/schema';
import { getCorrelationId } from '../common/correlation';
import { LOG_ENDPOINT } from '../common/handler';

type ClickstreamAction = ClickstreamLogEvent['details']['action'];
type ClickstreamTarget = ClickstreamLogEvent['details']['target'];

const SERVICE = import.meta.env.VITE_APP_NAME    || 'komodo-ui';
const VERSION = import.meta.env.VITE_APP_VERSION || 'unknown';
const ENV     = import.meta.env.MODE             || 'development';

export class ClickstreamLogger extends LoggingAdapter {
  constructor() {
    super({
      provider:  'clickstream',
      limit:     20,           // higher volume — batch more before flushing
      interval:  5_000,
      endpoint:  LOG_ENDPOINT,
    });
  }

  track(action: ClickstreamAction, target: ClickstreamTarget, requestId?: string): void {
    if (!browser || ENV === 'mock') return;

    const event: ClickstreamLogEvent = {
      timestamp:     new Date().toISOString(),
      level:         'info',
      type:          'clickstream',
      service:       SERVICE,
      env:           ENV,
      version:       VERSION,
      requestId,
      correlationId: getCorrelationId(),
      message:       `${action} on ${target.label || target.id || target.path || 'element'}`,
      details: {
        action,
        target,
        url:      window.location.href,
        viewport: `${window.innerWidth}x${window.innerHeight}`,
      },
    };

    this.send(event);
  }
}
