// Standard log schema shared across all Komodo frontend loggers.
// The JSON shape is what gets shipped to the server (and ultimately CloudWatch).
// Keep this in sync with the Go SDK's standard field set.

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';
export type LogEventType = 'runtime' | 'clickstream' | 'interaction' | 'telemetry';

// Base fields present on every log event regardless of type.
export interface BaseLogEvent {
  timestamp: string;        // ISO 8601 — set at emit time
  level: LogLevel;
  type: LogEventType;
  service: string;          // e.g. 'komodo-ui'
  env: string;              // 'production' | 'staging' | 'development' | 'mock'
  version: string;          // app version from VITE_APP_VERSION
  requestId?: string;       // E2E per HTTP request (from X-Request-ID response header)
  correlationId: string;    // browser fingerprint — sessionStorage UUID, resets on browser close
  userId?: string;          // populated when user is authenticated
  sessionId?: string;       // JWT session identifier
  message: string;
  details?: Record<string, unknown>;
}

// Runtime — service errors, unhandled exceptions, lifecycle events
export interface RuntimeLogEvent extends BaseLogEvent {
  type: 'runtime';
  details?: {
    error?: string;
    stack?: string;
    component?: string;
    request?: { method: string; path: string; status?: number };
    [key: string]: unknown;
  };
}

// Clickstream — raw DOM events (high-frequency, fine-grained UX analysis)
export interface ClickstreamLogEvent extends BaseLogEvent {
  type: 'clickstream';
  level: 'info'; // clickstream is always info
  details: {
    action: 'click' | 'hover' | 'scroll' | 'submit' | 'input' | 'focus' | 'blur';
    target: {
      id?: string;
      label?: string;
      text?: string;
      path?: string;      // DOM selector path
      aria?: string;
    };
    url: string;
    viewport?: string;    // e.g. '1440x900'
  };
}

// Interaction — semantic business events (lower-frequency, conversion tracking)
export interface InteractionLogEvent extends BaseLogEvent {
  type: 'interaction';
  level: 'info';
  details: {
    action: string;       // e.g. 'add_to_cart', 'checkout_start', 'search_submitted'
    url: string;
    data?: Record<string, unknown>; // event-specific payload (itemId, query, etc.)
  };
}

// Telemetry — performance metrics, traces, Core Web Vitals
export interface TelemetryLogEvent extends BaseLogEvent {
  type: 'telemetry';
  level: 'info';
  details: {
    name: string;         // metric or span name
    duration?: number;    // ms
    component?: string;
    traceId?: string;
    spanId?: string;
    [key: string]: unknown;
  };
}
