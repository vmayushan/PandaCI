import type { Method } from 'api-typify';

export interface GitRepository {
	id: string;
	name: string;
	description: string;
	public: boolean;
	url: string;
	updatedAt: string;
}

export interface GitRepositories {
	repos: GitRepository[];
	limit: number;
	limitExceeded: boolean;
}

export interface GitInstallation {
	id: string;
	name: string;
	avatarURL: string;
	isUser: boolean;
	repositoryScopes: 'all' | 'selected';
	type: 'github';
}

export interface GitInstallations {
	installations: GitInstallation[];
	isLastPage: boolean;
}

type GET = Method<{
	'/v1/git/github/installations': {
		res: GitInstallations;
		queries?: {
			page?: number;
			perPage?: number;
		};
	};
	'/v1/git/github/installations/{installationID}/repos': {
		res: GitRepositories;
		queries?: {
			query?: string;
			name?: string;
			owner?: string;
		};
	};
}>;

export interface GitAPI {
	GET: GET;
}
