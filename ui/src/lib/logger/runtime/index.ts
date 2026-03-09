import { SplunkLogger } from '../splunk';
import { OtelLogger } from '../otel';
import { LoggingAdapter } from '../common/adapter';

export type LoggerLevel = 'off' | 'debug' | 'info' | 'warn' | 'error' | 'panic';

const LoggerLevelWeight: Record<LoggerLevel, number> = { 
  off: 0, debug: 1, info: 2, warn: 3, error: 4, panic: 5 
};

type LoggerConfig = {
  level?: LoggerLevel;
  remoteLogger?: LoggingAdapter;
}

let levelWeight: number = (() => {
  const level = String(import.meta.env.LOG_LEVEL || 'error').toLowerCase() as LoggerLevel;
  return LoggerLevelWeight[level] ?? LoggerLevelWeight.error;
})();

let remoteLogger: LoggingAdapter | null = null;

const enableFullConsole = import.meta.env.DEV 
  || import.meta.env.MODE === 'development'
  || import.meta.env.MODE === 'qa'
  || import.meta.env.MODE === 'uat'
  || import.meta.env.MODE === 'staging';

const enableRemoteLogging = import.meta.env.PROD
  || import.meta.env.MODE === 'production'
  || import.meta.env.MODE === 'qa'
  || import.meta.env.MODE === 'uat'
  || import.meta.env.MODE === 'staging';

const logger = {
  init(config: LoggerConfig) {
    if (config.level) levelWeight = LoggerLevelWeight[config.level];
    
    if (enableRemoteLogging && config.remoteLogger) {
      if (config.remoteLogger instanceof LoggingAdapter) {
        remoteLogger = config.remoteLogger;
      } else if (typeof config.remoteLogger === 'string') {
        switch (config.remoteLogger) {
          case 'splunk':
            remoteLogger = new SplunkLogger();
            break;
          case 'otel':
            remoteLogger = new OtelLogger();
            break;
          default:
            throw new Error('Invalid remote logger');
        }
      }
    }
  },

  debug(message: string, ...args: any[]) {
    if (levelWeight >= LoggerLevelWeight.debug) {
      if (enableFullConsole) console.debug(message, ...args);
    }
  },

  trace(message: string, ...args: any[]) {
    if (levelWeight >= LoggerLevelWeight.debug) {
      if (enableFullConsole) console.trace(message, ...args);
    }
  },
  
  info(message: string, ...args: any[]) {
    if (levelWeight >= LoggerLevelWeight.info) {
      if (enableFullConsole) console.info(message, ...args);
      if (enableRemoteLogging && remoteLogger) remoteLogger.send({ level: 'info', message });
    }
  },
  
  warn(message: string, ...args: any[]) {
    if (levelWeight >= LoggerLevelWeight.warn) {
      if (enableFullConsole) console.warn(message, ...args);
      if (enableRemoteLogging && remoteLogger) remoteLogger.send({ level: 'warn', message });
    }
  },
  
  error(message: string, context?: Error | object) {
    if (levelWeight >= LoggerLevelWeight.error) {
      console.error(message, context);
      if (enableRemoteLogging && remoteLogger) { 
        remoteLogger.send({ level: 'error', message, context });
      }
    }
  },

  panic(message: string, context?: Error | object) {
    if (levelWeight >= LoggerLevelWeight.panic) {
      console.error(message, context);
      if (enableRemoteLogging && remoteLogger) {
        remoteLogger.send({ level: 'panic', message, context });
      }
    }
  },

  log(message: string, ...args: any[]) {
    console.log(message, ...args);
  }
};

export default logger;
