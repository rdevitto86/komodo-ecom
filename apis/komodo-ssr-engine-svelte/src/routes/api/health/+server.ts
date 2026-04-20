import { json } from '@sveltejs/kit';
import { S3Client, HeadBucketCommand } from '@aws-sdk/client-s3';
import { AWS_REGION, S3_CONTENT_BUCKET } from '$lib/config';

export async function GET() {
  const checks = {
    status: 'ok',
    timestamp: new Date().toISOString(),
    s3: 'unknown'
  };
  
  try {
    await new S3Client({ region: process.env[AWS_REGION] })
      .send(new HeadBucketCommand({ Bucket: process.env[S3_CONTENT_BUCKET]! }));
  } catch {
    checks.s3 = 'error';
    checks.status = 'degraded';
  }
  
  return json(checks, { status: checks.status === 'ok' ? 200 : 503 });
}