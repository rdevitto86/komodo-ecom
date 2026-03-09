import { S3Client, GetObjectCommand } from "@aws-sdk/client-s3";
import { logger } from '../logger';
import { getCachedTemplate, setCachedTemplate } from './cache';

const BUCKET = process.env.S3_CONTENT_BUCKET!;

const s3 = new S3Client({
  region: process.env.AWS_REGION || "us-east-1"
});

export async function getPageContentFromS3(pageKey: string) {
  // L1 (LRU) → L2 (SQLite read-only) — handled inside getCachedTemplate
  const cached = getCachedTemplate(pageKey);
  if (cached) return cached.content;

  // L3: Fetch from S3
  logger.info(`[S3 Fetch] ${pageKey}`);

  try {
    const command = new GetObjectCommand({
      Bucket: BUCKET,
      Key: `pages/${pageKey}.json`
    });

    const response = await s3.send(command);
    const bodyString = await response.Body?.transformToString();

    if (!bodyString) {
      throw new Error(`Empty response for ${pageKey}`);
    }

    const content = JSON.parse(bodyString);

    // Store in SQLite cache for subsequent requests
    setCachedTemplate(pageKey, content);

    return content;
  } catch (error) {
    logger.error(`Failed to fetch ${pageKey} from S3:`, error as Error);
    throw error;
  }
}

export async function fetchAndCacheFromS3(pageKey: string, isPreloaded: boolean = false): Promise<void> {
  try {
    const command = new GetObjectCommand({
      Bucket: BUCKET,
      Key: `pages/${pageKey}.json`
    });

    const response = await s3.send(command);
    const bodyString = await response.Body?.transformToString();

    if (!bodyString) {
      logger.error(`[Warmup] Empty response for ${pageKey}`);
      return;
    }

    const content = JSON.parse(bodyString);
    setCachedTemplate(pageKey, content, { isPreloaded, ttlMs: isPreloaded ? 24 * 60 * 60 * 1000 : undefined });
  } catch (error) {
    logger.error(`[Warmup] Failed to fetch ${pageKey}:`, error as Error);
  }
}
