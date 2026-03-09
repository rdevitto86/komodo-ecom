import { z } from 'zod';

const requestScehma = z.object({
  method: z.enum(['GET', 'POST', 'PUT', 'DELETE', 'PATCH']),
  body: z.union([z.string(), z.instanceof(Blob),  z.instanceof(ArrayBuffer), z.null()]).optional(),
  headers: z.record(z.string(), z.string()).optional(),
});

export function isRequestValid(url: string | URL, req: RequestInit): boolean {
  if (!(z.string().url().safeParse(url instanceof URL ? url.href : url)).success) return false;
  if (!requestScehma.safeParse(req).success) return false;
  return true;
}
