import type { LoggerLevel } from '../runtime';
import { LoggingAdapter } from '../common/adapter';

interface SplunkPayload {
  sourcetype: string;
  app: {
    name: string;
    version: string;
  };
  event: {
    message: string;
    severity: LoggerLevel;
    env: string;
    sessionId?: string;
    url: string;
    timestamp: number;
    details?: any;
  };
}

export class SplunkLogger extends LoggingAdapter {
  #sourcetype = import.meta.env.SPLUNK_SOURCETYPE || import.meta.env.APP_NAME || 'unknown';
  #app = { name: import.meta.env.APP_NAME || 'unknown', version: import.meta.env.VERSION || 'unknown' };

  constructor() {
    super({ 
      provider: 'splunk', 
      limit: import.meta.env.SPLUNK_LIMIT || 10,
      interval: import.meta.env.SPLUNK_INTERVAL || 5000,
      endpoint: import.meta.env.SPLUNK_ENDPOINT
    });
  }

  message(payload: SplunkPayload) {
    this.send({
      sourcetype: this.#sourcetype,
      app: this.#app,
      event: payload,
    });
  }
}
