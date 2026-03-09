import { getPageContentFromS3 } from '$lib/server/s3';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ setHeaders }) => {
  setHeaders({
    'cache-control': 'public, max-age=300, s-maxage=600, stale-while-revalidate=86400'
  });

  const content = await getPageContentFromS3('products');
  
  return { content };
};
