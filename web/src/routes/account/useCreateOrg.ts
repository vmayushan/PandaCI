import { goto } from '$app/navigation';
import { API } from '$lib/api';
import { queries } from '$lib/queries';
import { QueryClient, createMutation } from '@tanstack/svelte-query';
import type { CreateOrgFormData } from '../welcome/createOrgForm.svelte';

export function useCreateOrgMutation() {
	const queryClient = new QueryClient();

	return createMutation(() => ({
		mutationFn: (data: CreateOrgFormData) => API.post('/v1/orgs', { body: data }),
		onSuccess: async (org) => {
			queryClient.invalidateQueries(queries.organization.list());
			queryClient.setQueryData(queries.organization.getByName(org.slug).queryKey, org);
			await goto(`/${org.slug}`);
		}
	}));
}
