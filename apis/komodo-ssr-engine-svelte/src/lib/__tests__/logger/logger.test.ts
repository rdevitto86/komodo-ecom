import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

describe('logger', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    vi.resetModules();
    process.env = { ...originalEnv };
    vi.spyOn(console, 'log').mockImplementation(() => {});
    vi.spyOn(console, 'warn').mockImplementation(() => {});
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    process.env = originalEnv;
    vi.restoreAllMocks();
  });

  describe('info', () => {
    it('should log info messages when LOG_LEVEL is info', async () => {
      process.env.LOG_LEVEL = 'info';
      
      // Re-import to get fresh logger with new env
      const { logger } = await import('$lib/logger');
      
      logger.info('test message', { userId: '123' });
      
      expect(console.log).toHaveBeenCalledWith(
        JSON.stringify({
          level: 'info',
          msg: 'test message',
          userId: '123',
          timestamp: expect.any(Number)
        })
      );
    });

    it('should log info messages when LOG_LEVEL is warn', async () => {
      process.env.LOG_LEVEL = 'warn';
      
      const { logger } = await import('$lib/logger');
      
      logger.info('test message');
      
      expect(console.log).not.toHaveBeenCalled();
    });

    it('should not log info messages when LOG_LEVEL is error', async () => {
      process.env.LOG_LEVEL = 'error';
      
      const { logger } = await import('$lib/logger');
      
      logger.info('test message');
      
      expect(console.log).not.toHaveBeenCalled();
    });

    it('should not log info messages when LOG_LEVEL is off', async () => {
      process.env.LOG_LEVEL = 'off';
      
      const { logger } = await import('$lib/logger');
      
      logger.info('test message');
      
      expect(console.log).not.toHaveBeenCalled();
    });

    it('should handle info messages without metadata', async () => {
      process.env.LOG_LEVEL = 'info';
      
      const { logger } = await import('$lib/logger');
      
      logger.info('test message');
      
      expect(console.log).toHaveBeenCalledWith(
        JSON.stringify({
          level: 'info',
          msg: 'test message',
          timestamp: expect.any(Number)
        })
      );
    });
  });

  describe('warn', () => {
    it('should log warn messages when LOG_LEVEL is warn', async () => {
      process.env.LOG_LEVEL = 'warn';
      
      const { logger } = await import('$lib/logger');
      
      logger.warn('warning message', { userId: '123' });
      
      expect(console.warn).toHaveBeenCalledWith(
        JSON.stringify({
          level: 'warn',
          msg: 'warning message',
          userId: '123',
          timestamp: expect.any(Number)
        })
      );
    });

    it('should log warn messages when LOG_LEVEL is error', async () => {
      process.env.LOG_LEVEL = 'error';
      
      const { logger } = await import('$lib/logger');
      
      logger.warn('warning message');
      
      expect(console.warn).not.toHaveBeenCalled();
    });

    it('should not log warn messages when LOG_LEVEL is off', async () => {
      process.env.LOG_LEVEL = 'off';
      
      const { logger } = await import('$lib/logger');
      
      logger.warn('warning message');
      
      expect(console.warn).not.toHaveBeenCalled();
    });
  });

  describe('error', () => {
    it('should always log error messages regardless of LOG_LEVEL', async () => {
      process.env.LOG_LEVEL = 'off';
      
      const { logger } = await import('$lib/logger');
      
      const testError = new Error('test error');
      logger.error('error message', testError);
      
      expect(console.error).toHaveBeenCalledWith(
        JSON.stringify({
          level: 'error',
          msg: 'error message',
          error: 'test error',
          timestamp: expect.any(Number)
        })
      );
    });

    it('should log error messages without error object', async () => {
      const { logger } = await import('$lib/logger');
      
      logger.error('error message');
      
      expect(console.error).toHaveBeenCalledWith(
        JSON.stringify({
          level: 'error',
          msg: 'error message',
          error: undefined,
          timestamp: expect.any(Number)
        })
      );
    });

    it('should handle error with undefined error', async () => {
      const { logger } = await import('$lib/logger');
      
      logger.error('error message', undefined);
      
      expect(console.error).toHaveBeenCalledWith(
        JSON.stringify({
          level: 'error',
          msg: 'error message',
          error: undefined,
          timestamp: expect.any(Number)
        })
      );
    });
  });

  describe('default behavior', () => {
    it('should default to error level when LOG_LEVEL is not set', async () => {
      delete process.env.LOG_LEVEL;
      
      const { logger } = await import('$lib/logger');
      
      logger.info('info message');
      logger.warn('warn message');
      logger.error('error message');
      
      expect(console.log).not.toHaveBeenCalled();
      expect(console.warn).not.toHaveBeenCalled();
      expect(console.error).toHaveBeenCalled();
    });

    it('should default to error level when LOG_LEVEL is invalid', async () => {
      process.env.LOG_LEVEL = 'invalid';
      
      const { logger } = await import('$lib/logger');
      
      logger.info('info message');
      logger.warn('warn message');
      logger.error('error message');
      
      expect(console.log).not.toHaveBeenCalled();
      expect(console.warn).not.toHaveBeenCalled();
      expect(console.error).toHaveBeenCalled();
    });
  });
});