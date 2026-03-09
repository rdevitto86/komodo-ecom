import Database from 'better-sqlite3';
import { logger } from '$lib/logger';
import path from 'path';

export const DB_PATH = process.env.SQLITE_DB_PATH || path.join(process.cwd(), 'data', 'template-cache.db');

// ── Read-only connection for the main thread ────────────────────────────────
// The writer worker thread owns the read-write connection.
// WAL mode allows concurrent readers + single writer without blocking.

let readDb: Database.Database;

export function getReadDb(): Database.Database {
  if (!readDb) {
    readDb = new Database(DB_PATH, { readonly: true });
    readDb.pragma('journal_mode = WAL');
    readDb.pragma('cache_size = -64000');
    readDb.pragma('busy_timeout = 5000');
    readDb.pragma('temp_store = MEMORY');
    readDb.pragma('mmap_size = 268435456');
    logger.info(`[SQLite] Read-only connection opened at ${DB_PATH}`);
  }
  return readDb;
}

// ── Prepared read statements (created once, reused for every request) ───────

let _stmtReadTemplate: Database.Statement | null = null;

export function getReadTemplateStmt(): Database.Statement {
  if (!_stmtReadTemplate) {
    const db = getReadDb();
    _stmtReadTemplate = db.prepare(`
      SELECT t.id, t.content, t.content_hash, t.expires_at,
             cs.hit_count, cs.is_preloaded
      FROM templates t
      LEFT JOIN cache_stats cs ON cs.template_id = t.id
      WHERE t.route_key = ?
    `);
  }
  return _stmtReadTemplate;
}

export function closeReadDb(): void {
  if (readDb) {
    readDb.close();
    logger.info('[SQLite] Read-only connection closed');
  }
}
