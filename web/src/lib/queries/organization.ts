import { API } from '$lib/api';
import { createQueryKeys } from '@lukemorales/query-key-factory';

export const organizationQueries = createQueryKeys('organization', {
	list: () => ({
		queryFn: () => API.get('/v1/orgs', {}),
		queryKey: ['list']
	}),
	getByName: (orgSlug: string) => ({
		queryKey: [orgSlug],
		queryFn: () => API.get('/v1/orgs/{orgSlug}', { params: { orgSlug } }),
		contextQueries: {
			usage: () => ({
				queryKey: ['usage'],
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/usage', {
						params: { orgSlug }
					})
			}),
			users: () => ({
				queryKey: ['users'],
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/users', {
						params: { orgSlug }
					})
			}),
			projects: () => ({
				queryKey: ['projects'],
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects', {
						params: { orgSlug }
					})
			}),
			projectByName: (projectSlug: string) => ({
				queryKey: ['projects', 'name', projectSlug],
				queryFn: () =>
					API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}', {
						params: { orgSlug, projectSlug }
					})
			})
		}
	})
});
