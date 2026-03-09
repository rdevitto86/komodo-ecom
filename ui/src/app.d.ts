import 'vitest';
import type { Logger } from '$lib/logger/splunk';
declare global {
  namespace App {
    interface Locals {
      secrets: Record<string, string>;
    }
    // interface Error {}
    // interface PageData {}
    // interface Platform {}
  }
  var logger: Logger;
}

declare module '$env/dynamic/private' {
  export const STRIPE_KEY: string;
  export const BACKEND_TOKEN: string;
}

declare module '$env/dynamic/public' {
  export const PUBLIC_API_URL: string;
}

declare module '*.css' {
  const classes: { [key: string]: string };
  export default classes;
}

declare module '*.png' {
  const value: string;
  export default value;
}

declare module 'vitest' {
  export interface TaskMeta {
    tags?: string[];
    type?: 'integration' | 'component' | 'unit';
    requirementId?: string;
  }
}

export {};