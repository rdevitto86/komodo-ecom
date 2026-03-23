import 'vitest';
declare global {
  namespace App {
    interface Locals {
      requestId: string;
    }
    // interface Error {}
    // interface PageData {}
    // interface Platform {}
  }
}

declare module '$env/dynamic/private' {
  export const AUTH_API_URL: string;
  export const USER_API_URL: string;
  export const SHOP_ITEMS_API_URL: string;
  export const CART_API_URL: string;
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