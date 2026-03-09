// TODO pino logger?

const LOG_LEVEL = process.env.LOG_LEVEL as keyof typeof LogLevel || 'error';
enum LogLevel { off = 0, info = 1, warn = 2, error = 3 }

export const logger = {
  info: (msg: string, meta?: object) => (LogLevel.info >= LogLevel[LOG_LEVEL])
    ? console.log(JSON.stringify({ level: 'info', msg, ...meta, timestamp: Date.now() })) : null,

  warn: (msg: string, meta?: object) => (LogLevel.warn >= LogLevel[LOG_LEVEL])
    ? console.warn(JSON.stringify({ level: 'warn', msg, ...meta, timestamp: Date.now() })) : null,

  error: (msg: string, error?: Error) => console.error(JSON.stringify({
    level: 'error',
    msg,
    error: error?.message,
    timestamp: Date.now()
  }))
};
