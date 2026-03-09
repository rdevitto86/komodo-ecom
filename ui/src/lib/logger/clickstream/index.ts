import { LoggingAdapter } from '../common/adapter';

interface ClickstreamPayload {
  action: 'click' | 'hover' | 'scroll' | 'submit' | 'input';
  target: {
    id?: string;
    label?: string;
    text?: string;
    path?: string;
    aria?: string;
    [key: string]: any;
  };
  meta: {
    screenResolution?: string;
    viewport?: string;
    userLanguage?: string;
    theme?: 'light' | 'dark';
    timezone?: string;
    os?: string;
    osVersion?: string;
    browser?: string;
    browserVersion?: string;
    [key: string]: any;
  };
}

export class ClickstreamLogger extends LoggingAdapter {
  constructor() {
    super({ 
      provider: 'clickstream', 
      limit: 10, 
      interval: 5000,
      endpoint: '/api/logging/clickstream'
    });
  }

  trace(event: ClickstreamPayload): void {
    this.send(event);
  }
}
