import { browser } from '$app/environment';
import { LoggingAdapter } from '../common/adapter';
import type { InteractionLogEvent } from '../common/schema';
import { getCorrelationId } from '../common/correlation';
import { LOG_ENDPOINT } from '../common/handler';

const SERVICE = import.meta.env.VITE_APP_NAME    || 'komodo-ui';
const VERSION = import.meta.env.VITE_APP_VERSION || 'unknown';
const ENV     = import.meta.env.MODE             || 'development';

/**
 * InteractionLogger — semantic business events.
 *
 * Use this for user-intent actions, not raw DOM events (use ClickstreamLogger
 * for those). Examples: 'add_to_cart', 'checkout_start', 'search_submitted',
 * 'coupon_applied', 'review_submitted'.
 */
export class InteractionLogger extends LoggingAdapter {
  constructor() {
    super({
      provider:  'interaction',
      limit:     10,
      interval:  10_000,
      endpoint:  LOG_ENDPOINT,
    });
  }

  track(action: string, data?: Record<string, unknown>, requestId?: string): void {
    if (!browser || ENV === 'mock') return;

    const event: InteractionLogEvent = {
      timestamp:     new Date().toISOString(),
      level:         'info',
      type:          'interaction',
      service:       SERVICE,
      env:           ENV,
      version:       VERSION,
      requestId,
      correlationId: getCorrelationId(),
      message:       action,
      details: {
        action,
        url:  window.location.href,
        data,
      },
    };

    this.send(event);
  }
}
