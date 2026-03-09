import { parentPort } from 'node:worker_threads';
import Database from 'better-sqlite3';
import path from 'path';

/**
 * Cache writer worker thread
 * 
 * This worker thread is responsible for writing to the SQLite database.
 * It is a separate thread to avoid blocking the main thread.
 */

if (!parentPort) throw new Error('cache-writer must run as a worker thread');

const DB_PATH = process.env.SQLITE_DB_PATH || path.join(process.cwd(), 'data', 'template-cache.db');
const HIT_FLUSH_INTERVAL_MS = 5_000;
const MAINTENANCE_INTERVAL_MS = 10 * 60 * 1000;
const PROMOTION_THRESHOLD = 5;

const db = new Database(DB_PATH);

db.pragma('journal_mode = WAL');
db.pragma('synchronous = NORMAL');
db.pragma('cache_size = -64000');
db.pragma('busy_timeout = 5000');
db.pragma('temp_store = MEMORY');
db.pragma('mmap_size = 268435456');
db.pragma('foreign_keys = ON');

// ── Build DB schema if it doesn't exist ────────────────────────────────────

db.exec(`
  CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    route_key TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at TEXT NOT NULL
  );

  CREATE INDEX IF NOT EXISTS idx_templates_route_key ON templates(route_key);
  CREATE INDEX IF NOT EXISTS idx_templates_expires_at ON templates(expires_at);

  CREATE TABLE IF NOT EXISTS cache_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    template_id INTEGER NOT NULL UNIQUE,
    hit_count INTEGER NOT NULL DEFAULT 0,
    last_accessed_at TEXT NOT NULL DEFAULT (datetime('now')),
    is_preloaded INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE
  );

  CREATE INDEX IF NOT EXISTS idx_cache_stats_hit_count ON cache_stats(hit_count DESC);
  CREATE INDEX IF NOT EXISTS idx_cache_stats_template_id ON cache_stats(template_id);
`);

// ── Prepared statements (created once, reused forever) ──────────────────────

const stmtUpsertTemplate = db.prepare(`
  INSERT INTO templates (route_key, content, content_hash, updated_at, expires_at)
  VALUES (?, ?, ?, datetime('now'), ?)
  ON CONFLICT(route_key) DO UPDATE SET
    content = excluded.content,
    content_hash = excluded.content_hash,
    updated_at = datetime('now'),
    expires_at = excluded.expires_at
`);

const stmtGetTemplateId = db.prepare(
  'SELECT id FROM templates WHERE route_key = ?'
);

const stmtUpsertStats = db.prepare(`
  INSERT INTO cache_stats (template_id, hit_count, last_accessed_at, is_preloaded)
  VALUES (?, 0, datetime('now'), ?)
  ON CONFLICT(template_id) DO UPDATE SET
    is_preloaded = CASE WHEN excluded.is_preloaded = 1 THEN 1 ELSE cache_stats.is_preloaded END,
    last_accessed_at = datetime('now')
`);

const stmtBatchHitUpdate = db.prepare(`
  UPDATE cache_stats
  SET hit_count = hit_count + ?, last_accessed_at = datetime('now')
  WHERE template_id = (SELECT id FROM templates WHERE route_key = ?)
`);

const stmtPurgeExpired = db.prepare(`
  DELETE FROM templates
  WHERE expires_at <= datetime('now')
  AND id NOT IN (
    SELECT template_id FROM cache_stats WHERE is_preloaded = 1
  )
`);

const stmtFrequentRoutes = db.prepare(`
  SELECT t.route_key
  FROM cache_stats cs
  JOIN templates t ON t.id = cs.template_id
  WHERE cs.hit_count >= ? AND cs.is_preloaded = 0
  ORDER BY cs.hit_count DESC
  LIMIT ?
`);

const stmtPromote = db.prepare(
  'UPDATE cache_stats SET is_preloaded = 1 WHERE template_id = (SELECT id FROM templates WHERE route_key = ?)'
);

const stmtExtendTTL = db.prepare(
  "UPDATE templates SET expires_at = datetime('now', '+24 hours') WHERE route_key = ?"
);

const stmtCountTotal = db.prepare('SELECT COUNT(*) as count FROM templates');
const stmtCountPreloaded = db.prepare('SELECT COUNT(*) as count FROM cache_stats WHERE is_preloaded = 1');
const stmtCountFrequent = db.prepare(`SELECT COUNT(*) as count FROM cache_stats WHERE hit_count >= ${PROMOTION_THRESHOLD}`);
const stmtCountExpired = db.prepare("SELECT COUNT(*) as count FROM templates WHERE expires_at <= datetime('now')");

// ── Transactions ────────────────────────────────────────────────────────────

const txUpsertTemplate = db.transaction(
  (routeKey: string, serialized: string, hash: string, expires: string, isPreloaded: boolean) => {
    stmtUpsertTemplate.run(routeKey, serialized, hash, expires);
    const row = stmtGetTemplateId.get(routeKey) as { id: number };
    stmtUpsertStats.run(row.id, isPreloaded ? 1 : 0);
  }
);

const txFlushHits = db.transaction((hits: Map<string, number>) => {
  for (const [routeKey, count] of hits) {
    stmtBatchHitUpdate.run(count, routeKey);
  }
});

const txMaintenance = db.transaction(() => {
  const purged = stmtPurgeExpired.run();

  const frequent = stmtFrequentRoutes.all(PROMOTION_THRESHOLD, 20) as { route_key: string }[];
  for (const row of frequent) {
    stmtPromote.run(row.route_key);
    stmtExtendTTL.run(row.route_key);
  }

  return { purged: purged.changes, promoted: frequent.length };
});

// ── Hit count accumulator (flushed periodically) ────────────────────────────

const pendingHits = new Map<string, number>();

function recordHit(routeKey: string): void {
  pendingHits.set(routeKey, (pendingHits.get(routeKey) || 0) + 1);
}

function flushHits(): void {
  if (pendingHits.size === 0) return;
  const batch = new Map(pendingHits);
  pendingHits.clear();
  txFlushHits(batch);
}

// ── Periodic timers ─────────────────────────────────────────────────────────

setInterval(flushHits, HIT_FLUSH_INTERVAL_MS);
setInterval(() => {
  try {
    flushHits(); // flush before maintenance so counts are current
    const result = txMaintenance();
    parentPort!.postMessage({ type: 'maintenance-done', ...result });
  } catch (err) {
    parentPort!.postMessage({ type: 'error', message: (err as Error).message });
  }
}, MAINTENANCE_INTERVAL_MS);

// ── Message handler ─────────────────────────────────────────────────────────

export type WriterMessage =
  | { type: 'upsert'; routeKey: string; content: string; hash: string; expires: string; isPreloaded: boolean }
  | { type: 'hit'; routeKey: string }
  | { type: 'flush-hits' }
  | { type: 'maintenance' }
  | { type: 'get-stats'; requestId: string }
  | { type: 'get-frequent'; requestId: string; limit: number }
  | { type: 'promote'; routeKey: string }
  | { type: 'shutdown' };

parentPort.on('message', (msg: WriterMessage) => {
  try {
    switch (msg.type) {
      case 'upsert':
        txUpsertTemplate(msg.routeKey, msg.content, msg.hash, msg.expires, msg.isPreloaded);
        break;

      case 'hit':
        recordHit(msg.routeKey);
        break;

      case 'flush-hits':
        flushHits();
        break;

      case 'maintenance': {
        flushHits();
        const result = txMaintenance();
        parentPort!.postMessage({ type: 'maintenance-done', ...result });
        break;
      }

      case 'get-stats': {
        const total = (stmtCountTotal.get() as { count: number }).count;
        const preloaded = (stmtCountPreloaded.get() as { count: number }).count;
        const frequent = (stmtCountFrequent.get() as { count: number }).count;
        const expired = (stmtCountExpired.get() as { count: number }).count;
        parentPort!.postMessage({
          type: 'stats',
          requestId: msg.requestId,
          stats: { total, preloaded, frequent, expired },
        });
        break;
      }

      case 'get-frequent': {
        const rows = stmtFrequentRoutes.all(PROMOTION_THRESHOLD, msg.limit) as { route_key: string }[];
        parentPort!.postMessage({
          type: 'frequent-routes',
          requestId: msg.requestId,
          routes: rows.map(r => r.route_key),
        });
        break;
      }

      case 'promote':
        stmtPromote.run(msg.routeKey);
        stmtExtendTTL.run(msg.routeKey);
        break;

      case 'shutdown':
        flushHits();
        db.close();
        process.exit(0);
    }
  } catch (err) {
    parentPort!.postMessage({ type: 'error', message: (err as Error).message });
  }
});

parentPort.postMessage({ type: 'ready' });
