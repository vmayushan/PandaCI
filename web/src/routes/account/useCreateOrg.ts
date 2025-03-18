import { goto } from '$app/navigation';
import { API, type Organization } from '$lib/api';
import { queries } from '$lib/queries';
import { QueryClient, createMutation } from '@tanstack/svelte-query';
import type { CreateOrgFormData } from '../welcome/createOrgForm.svelte';

export function useCreateOrgMutation() {
	const queryClient = new QueryClient();

	return createMutation(() => ({
		mutationFn: (data: CreateOrgFormData) => API.post('/v1/orgs', { body: data }),
		onSuccess: (org) => {
			queryClient.setQueryData(queries.organization.getByName(org.slug).queryKey, org);
			queryClient.setQueryData(queries.organization.list().queryKey, (old: Organization[] = []) => [
				...old,
				org
			]);

			return goto(`/${org.slug}/new`);
		},
		onSettled: () => {
			queryClient.invalidateQueries(queries.organization.list());
		}
	}));
}
