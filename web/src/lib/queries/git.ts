import { API, type apiDefs } from '$lib/api';
import { createQueryKeys } from '@lukemorales/query-key-factory';
import { isEmpty } from 'lodash-es';

export const githubQueries = createQueryKeys('github', {
	listInstallation: (queries: apiDefs['GET']['/v1/git/github/installations']['queries'] = {}) => ({
		queryFn: () =>
			API.get('/v1/git/github/installations', { queries }) as Promise<
				apiDefs['GET']['/v1/git/github/installations']['res']
			>,
		queryKey: [
			'list',
			...(isEmpty(queries) ? [] : [JSON.stringify(queries, Object.keys(queries).sort())])
		],
		contextQueries: {
			listRepositories: (
				installationID: string,
				queries: apiDefs['GET']['/v1/git/github/installations/{installationID}/repos']['queries'] = {}
			) => ({
				queryFn: () =>
					API.get('/v1/git/github/installations/{installationID}/repos', {
						params: {
							installationID
						},
						queries
					}) as Promise<
						apiDefs['GET']['/v1/git/github/installations/{installationID}/repos']['res']
					>,
				queryKey: [
					installationID,
					...(isEmpty(queries) ? [] : [JSON.stringify(queries, Object.keys(queries).sort())])
				]
			})
		}
	})
});
