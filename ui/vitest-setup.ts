import { vi, beforeEach, beforeAll } from 'vitest';
import '@testing-library/jest-dom';
import { BASE_API_URL as BASE_API_URL_KEY } from './src/lib/config';

// globalThis.fetch = vi.fn();

beforeAll(() => {
  process.env[BASE_API_URL_KEY] = 'https://api.example.com';
});

beforeEach(() => {
  vi.clearAllMocks();
});
