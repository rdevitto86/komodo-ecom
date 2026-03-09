import { vi, beforeEach, beforeAll } from 'vitest';
import '@testing-library/jest-dom';

// globalThis.fetch = vi.fn();

beforeAll(() => {
  process.env.BASE_API_URL = 'https://api.example.com';
});

beforeEach(() => {
  vi.clearAllMocks();
});
