import { env } from '$env/dynamic/private';
import { APIClient } from './common/client';

interface SecretsResponse {
	ENV: string;
	PUBLIC_API_URL: string;
	[key: string]: string;
}

export class SecretsAPI extends APIClient {

	async getSecrets(): Promise<SecretsResponse> {
		try {
			// return await this.send('GET', '/secrets');
			const res = {
				ENV:  env.ENV || 'local',
				PUBLIC_API_URL: env.PUBLIC_API_URL || '',
				AWS_REGION: 'us-east-1',
				AWS_ACCESS_KEY_ID: 'test',
				AWS_SECRET_ACCESS_KEY: 'test'
			};
			return res;
		} catch (err) {
			console.error('failed to fetch secrets from AWS Secrets Manager:', err as Error);
			throw err;
		}
	}
}
