import { Worker } from 'node:worker_threads';
import { getReadTemplateStmt } from './db';
import { LRUCache } from './lru';
import { logger } from '$lib/logger';
import crypto from 'crypto';
import path from 'path';
import { fileURLToPath } from 'url';

/**
 * In-memory LRU cache implementation for SQLite templates
 */

const DEFAULT_TTL_MS = 30 * 60 * 1000; // 30 minutes
const PRELOADED_TTL_MS = 24 * 60 * 60 * 1000; // 24 hours
const PROMOTION_THRESHOLD = 5;
const L1_MAX_SIZE = 500;

interface TemplateReadRow {
  id: number;
  content: string;
  content_hash: string;
  expires_at: string;
  hit_count: number | null;
  is_preloaded: number | null;
}

export interface CachedTemplate {
  content: any;
  hash: string;
  isPreloaded: boolean;
  hitCount: number;
}

// ── L1: In-memory LRU (parsed objects, zero I/O) ───────────────────────────

const l1 = new LRUCache<CachedTemplate>(L1_MAX_SIZE, DEFAULT_TTL_MS);

// ── Writer worker thread ────────────────────────────────────────────────────

let writer: Worker | null = null;
const pendingRequests = new Map<string, { resolve: (value: any) => void; reject: (err: Error) => void }>();
let requestIdCounter = 0;

function getWriterPath(): string {
  const dir = path.dirname(fileURLToPath(import.meta.url));
  return path.join(dir, 'writer.ts');
}

export function initCacheWriter(): Promise<void> {
  return new Promise((resolve, reject) => {
    writer = new Worker(getWriterPath());

    writer.on('message', (msg: any) => {
      if (msg.type === 'ready') {
        logger.info('[CacheWriter] Worker thread ready');
        resolve();
        return;
      }

      if (msg.type === 'error') {
        logger.error(`[CacheWriter] ${msg.message}`);
        return;
      }

      if (msg.type === 'maintenance-done') {
        logger.info(`[CacheWriter] Maintenance: purged=${msg.purged}, promoted=${msg.promoted}`);
        return;
      }

      // Resolve pending request/response pairs (stats, frequent-routes)
      if (msg.requestId) {
        const pending = pendingRequests.get(msg.requestId);
        if (pending) {
          pendingRequests.delete(msg.requestId);
          pending.resolve(msg);
        }
      }
    });

    writer.on('error', (err: Error) => {
      logger.error('[CacheWriter] Worker error', err);
      reject(err);
    });

    writer.on('exit', (code) => {
      logger.warn(`[CacheWriter] Worker exited with code ${code}`);
      writer = null;
    });
  });
}

function postToWriter(msg: any): void {
  if (!writer) {
    logger.error('[CacheWriter] Worker not initialized, dropping message');
    return;
  }
  writer.postMessage(msg);
}

function requestFromWriter<T>(msg: any): Promise<T> {
  return new Promise((resolve, reject) => {
    const requestId = String(++requestIdCounter);
    msg.requestId = requestId;
    pendingRequests.set(requestId, { resolve, reject });
    postToWriter(msg);

    // Timeout to prevent leaked promises
    setTimeout(() => {
      if (pendingRequests.has(requestId)) {
        pendingRequests.delete(requestId);
        reject(new Error(`[CacheWriter] Request ${requestId} timed out`));
      }
    }, 5000);
  });
}

// ── Helpers ─────────────────────────────────────────────────────────────────

function contentHash(content: string): string {
  return crypto.createHash('sha256').update(content).digest('hex').slice(0, 16);
}

function expiresAt(ttlMs: number = DEFAULT_TTL_MS): string {
  return new Date(Date.now() + ttlMs).toISOString();
}

function isExpired(expiresAtStr: string): boolean {
  return new Date(expiresAtStr).getTime() <= Date.now();
}

// ── Public API ──────────────────────────────────────────────────────────────

/**
 * L1 (LRU Map) → L2 (SQLite read-only) → null
 * Writes (hit count bump) offloaded to worker — zero writes on the hot path.
 */
export function getCachedTemplate(routeKey: string): CachedTemplate | null {
  // L1: in-memory (already parsed, ~0.001ms)
  const l1Hit = l1.get(routeKey);
  if (l1Hit) {
    postToWriter({ type: 'hit', routeKey });
    return l1Hit;
  }

  // L2: SQLite read-only connection (~0.1ms)
  const row = getReadTemplateStmt().get(routeKey) as TemplateReadRow | undefined;
  if (!row) return null;

  if (isExpired(row.expires_at) && !row.is_preloaded) {
    return null;
  }

  const template: CachedTemplate = {
    content: JSON.parse(row.content),
    hash: row.content_hash,
    isPreloaded: !!row.is_preloaded,
    hitCount: (row.hit_count || 0) + 1,
  };

  // Populate L1 for subsequent requests
  const ttl = row.is_preloaded ? PRELOADED_TTL_MS : DEFAULT_TTL_MS;
  l1.set(routeKey, template, ttl);

  // Async hit count bump via worker (no blocking write)
  postToWriter({ type: 'hit', routeKey });

  // Auto-promote if past threshold
  if (!row.is_preloaded && (row.hit_count || 0) + 1 >= PROMOTION_THRESHOLD) {
    postToWriter({ type: 'promote', routeKey });
  }

  return template;
}

/**
 * Writes to L1 immediately (instant availability) and to SQLite via worker.
 */
export function setCachedTemplate(
  routeKey: string,
  content: any,
  options: { isPreloaded?: boolean; ttlMs?: number } = {}
) {
  const { isPreloaded = false, ttlMs = DEFAULT_TTL_MS } = options;
  const serialized = JSON.stringify(content);
  const hash = contentHash(serialized);
  const expires = expiresAt(ttlMs);

  // L1: instant availability for next request
  const template: CachedTemplate = {
    content,
    hash,
    isPreloaded,
    hitCount: 0,
  };
  l1.set(routeKey, template, ttlMs);

  // L2: durable write via worker (non-blocking)
  postToWriter({
    type: 'upsert',
    routeKey,
    content: serialized,
    hash,
    expires,
    isPreloaded,
  });

  logger.info(`[TemplateCache] SET: ${routeKey} (preloaded: ${isPreloaded})`);
}

/**
 * Request frequent routes from the worker (reads from its connection).
 */
export async function getFrequentRoutes(limit: number = 20): Promise<string[]> {
  try {
    const response = await requestFromWriter<{ routes: string[] }>({
      type: 'get-frequent',
      limit,
    });
    return response.routes;
  } catch {
    return [];
  }
}

/**
 * Trigger promotion via worker (write operation).
 */
export function promoteToPreloaded(routeKey: string) {
  postToWriter({ type: 'promote', routeKey });
  logger.info(`[TemplateCache] PROMOTED to preloaded: ${routeKey}`);
}

/**
 * Trigger maintenance via worker (purge + promote in a single transaction).
 */
export function runMaintenance() {
  postToWriter({ type: 'maintenance' });
}

/**
 * Request cache stats from the worker.
 */
export async function getCacheStats(): Promise<{ total: number; preloaded: number; frequent: number; expired: number }> {
  try {
    const response = await requestFromWriter<{ stats: { total: number; preloaded: number; frequent: number; expired: number } }>({
      type: 'get-stats',
    });
    return response.stats;
  } catch {
    return { total: 0, preloaded: 0, frequent: 0, expired: 0 };
  }
}

/**
 * Graceful shutdown: flush pending hits and close worker.
 */
export function shutdownCache() {
  postToWriter({ type: 'shutdown' });
}

/**
 * Clear L1 cache (e.g., after admin content invalidation).
 */
export function invalidateL1(routeKey?: string) {
  if (routeKey) {
    l1.delete(routeKey);
  } else {
    l1.clear();
  }
}
