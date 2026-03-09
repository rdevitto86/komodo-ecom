import { LoggingAdapter } from '../common/adapter';

interface OtelPayload {
  name: string;
  attributes: {
    component?: string;
    traceId?: string;
    spanId?: string;
    duration?: number;
    memoryUsage?: number;
    [key: string]: any;
  }
}

export class OtelLogger extends LoggingAdapter {
  constructor() {
    super({ 
      provider: 'otel', 
      limit: 10, 
      interval: 5000,
      endpoint: '/api/logging/otel'
    });
  }

  trace(name: string, attributes?: OtelPayload['attributes']): void {
    this.send({ name, attributes: attributes || {} });
  }
}
