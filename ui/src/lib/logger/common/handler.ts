import logger from '../';

export type LogProviderType = 'splunk' | 'otel' | 'clickstream';
export type WorkerDirectives = 'CONFIG' | 'CONFIG_ACK' | 'LOG' | 'LOG_ACK' | 'FLUSH' | 'FLUSH_ACK' | 'STOP' | 'ERROR';
export type WorkerMessagePayload = { provider: LogProviderType; data: any };

// TODO - rehydrate logs from localStorage

export class WorkerHandler {
  static #instance: WorkerHandler | null = null;
  static #worker: Worker | null = null;

  private constructor() {
    if (!WorkerHandler.#instance) {
      WorkerHandler.#initWorker();
      WorkerHandler.#instance = this;
    }
    return WorkerHandler.#instance;
  }

  static init() {
    return new WorkerHandler();
  }

  static config(provider: LogProviderType, config: { limit: number; interval: number; endpoint: string }) {
    if (!this.#worker) this.#initWorker();
    if (!this.#worker) return;
    this.#worker.postMessage({
      directive: 'CONFIG',
      payload: ({ provider, data: config } as WorkerMessagePayload)
    });
  }

  static start() {
    if (this.#worker) return;
    this.#initWorker();
  }

  static stop() {
    if (!this.#worker) return;
    this.#worker.postMessage({ directive: 'STOP' });
    logger.info('Log worker stopped');
    this.#worker = null;
  }

  static send(provider: LogProviderType, data: any) {
    if (!this.#worker) this.#initWorker();
    if (!this.#worker || !data) return;
    this.#worker.postMessage({ directive: 'LOG', payload: ({ provider, data } as WorkerMessagePayload) });
  }

  static flush(provider: LogProviderType) {
    if (!this.#worker) this.#initWorker();
    if (!this.#worker) return;
    this.#worker.postMessage({ directive: 'FLUSH', payload: ({ provider } as WorkerMessagePayload) });
  }

  static isActive() {
    return this.#worker !== null;
  }

  static #initWorker() {
    if (this.#worker) return this.#worker;
    this.#worker = new Worker(new URL('./handler.worker.ts', import.meta.url), { type: 'module' });

    this.#worker.onmessage = ({ data: { directive, provider, payload } }) => {
      switch (directive) {          
        case 'ERROR':
          logger.error(`${provider} logger error`, payload?.error || payload);
          break;
          
        case 'STOP':
          const lostCounts = payload || {};
          const lostProviders = Object.entries(lostCounts)
            .filter(([_, count]) => (count as number) > 0)
            .map(([p, count]) => `${p}: ${count}`);
          
          if (lostProviders.length > 0) {
            logger.warn(`Log worker stopped with unsent logs: ${lostProviders.join(', ')}`);
          } else {
            logger.info('Log worker stopped cleanly');
          }
          break;
      }
    };

    this.#worker.onerror = (error) => logger.error('Log worker error', error);
    
    return this.#worker;
  }
}
