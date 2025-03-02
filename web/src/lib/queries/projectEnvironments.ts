import { API } from '$lib/api';
import { createQueryKeys } from '@lukemorales/query-key-factory';

export const projectEnvironmentQueries = createQueryKeys('environments', {
	projectEnvironments: (orgSlug: string, projectSlug: string) => ({
		queryKey: [orgSlug, projectSlug],
		contextQueries: {
			list: () => ({
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/environments', {
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
