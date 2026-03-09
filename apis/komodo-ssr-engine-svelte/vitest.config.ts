import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: ['src/**/__tests__/**/*.{test,spec}.{js,ts}'],
    environment: 'node',
    globals: true,
  }
});
