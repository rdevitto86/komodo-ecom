import { WorkerHandler } from './handler';
import type { LogProviderType } from './handler';
import { redact } from './redaction';

export type LoggingAdapterConfig = {
  provider: LogProviderType;
  limit: number;
  interval: number;
  endpoint: string;
}

export abstract class LoggingAdapter {
  protected worker: WorkerHandler | null = null;
  protected provider: LogProviderType;
  protected limit: number;
  protected interval: number;
  protected endpoint: string;

  constructor(config: LoggingAdapterConfig) {
    this.provider = config.provider;
    this.limit = config.limit;
    this.interval = config.interval;
    this.endpoint = config.endpoint;
    this.worker = WorkerHandler.init();
    
    if (this.worker) {
      WorkerHandler.config(this.provider, {
        limit: this.limit,
        interval: this.interval,
        endpoint: this.endpoint
      });
    }
  }

  send(payload: any): void {
    WorkerHandler.send(this.provider, redact(payload));
  }
}
