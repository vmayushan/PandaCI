import { API } from '$lib/api';
import { createQueryKeys } from '@lukemorales/query-key-factory';

export const runQueries = createQueryKeys('runs', {
	projectRuns: (orgSlug: string, projectSlug: string) => ({
		queryKey: [orgSlug, projectSlug],
		contextQueries: {
			list: () => ({
				queryFn: ({ pageParam }: { pageParam: number }) =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/runs', {
						params: {
							orgSlug,
							projectSlug
						},
						queries: {
							page: pageParam
						}
					}),
				queryKey: ['list']
			}),
			logs: (workflowID: string, logID: string) => ({
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/run/{workflowID}/logs/{logID}', {
						params: {
							orgSlug,
							projectSlug,
							workflowID,
							logID
						}
					}),
				queryKey: [workflowID, logID]
			}),

			get: (runNumber: string) => ({
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/runs/{runNumber}', {
						params: {
							orgSlug,
							projectSlug,
							runNumber
						}
					}),
				queryKey: [runNumber]
			})
		}
	})
});
