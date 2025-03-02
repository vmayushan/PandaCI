import { API } from '$lib/api';
import { createQueryKeys } from '@lukemorales/query-key-factory';

export const projectVariableQueries = createQueryKeys('variables', {
	projectVariables: (orgSlug: string, projectSlug: string) => ({
		queryKey: [orgSlug, projectSlug],
		contextQueries: {
			get: (variableID: string) => ({
				queryKey: [variableID],
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/variables/{variableID}', {
						params: {
							orgSlug,
							projectSlug,
							variableID
						},
						queries: {
							decrypt: true
						}
					})
			}),
			list: () => ({
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/variables', {
						params: {
							orgSlug,
							projectSlug
						}
					}),
				queryKey: ['list']
			})
		}
	})
});
