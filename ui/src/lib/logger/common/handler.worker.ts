import type { LogProviderType } from './handler';

type ProviderState = {
  config: { limit: number; interval: number; endpoint: string } | null;
  buffer: any[];
  interval: number | null;
};

const providers: Map<LogProviderType, ProviderState> = new Map();

// TODO - rehydrate logs from localStorage

/**
 * @param {MessageEvent<any>} param0 
 */
self.onmessage = async ({ data: { directive, provider, payload } }) => {
  try {
    switch (directive) {
      case 'CONFIG':
        if (!provider) throw new Error('Provider is required');
        
        const state = providers.get(provider);
        if (state) {
          // Clear existing interval if reconfiguring
          if (state.interval) clearInterval(state.interval);
        }
        
        const newState: ProviderState = {
          config: payload,
          buffer: state?.buffer || [],
          interval: null
        };
        
        if (payload?.interval) {
          newState.interval = setInterval(() => flush(provider), payload.interval) as unknown as number;
        }
        
        providers.set(provider, newState);
        break;
        
      case 'LOG':
        if (!provider) throw new Error('Provider is required');
        const logState = providers.get(provider);
        if (!logState?.config) throw new Error(`Provider ${provider} not configured`);
        
        logState.buffer.push(payload);
        if (logState.buffer.length >= logState.config.limit) {
          await flush(provider);
        }
        break;
        
      case 'FLUSH':
        if (!provider) throw new Error('Provider is required');
        await flush(provider);
        break;
        
      case 'STOP':
        // Flush all providers before stopping
        if (payload === true) {
          await Promise.all(
            Array.from(providers.keys()).map(async (pvdr) => await flush(pvdr))
          );
        }

        const lostCounts: Record<string, number> = {};
        providers.forEach((state, pvdr) => {
          if (state.buffer.length > 0) {
            lostCounts[pvdr] = state.buffer.length;
          }
          if (state.interval) clearInterval(state.interval);
        });

        self.postMessage({ directive: 'STOP', payload: lostCounts });
        setTimeout(() => self.close(), 100);
        break;
    }
  } catch (err) {
    self.postMessage({ directive: 'ERROR', provider: payload?.provider, payload: err });
  }
};

async function flush(provider: LogProviderType) {
  const state = providers.get(provider);
  if (!state?.buffer.length || !state.config) return;
  
  const logs = [...state.buffer];
  state.buffer = [];

  try {
    const response = await fetch(state.config.endpoint, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(logs),
      keepalive: true
    });

    if (!response.ok) {
      throw new Error(`Flush failed: ${response.status} - ${response.statusText}`);
    }
  } catch (err) {
    state.buffer = [...logs, ...state.buffer];
    self.postMessage({ 
      directive: 'ERROR', 
      provider,
      payload: { error: err, count: state.buffer.length } 
    });
  }
}
