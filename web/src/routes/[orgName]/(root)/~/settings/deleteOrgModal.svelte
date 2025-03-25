<script lang="ts">
	import {
		Dialog,
		DialogTitle,
		DialogActions,
		Text,
		Button,
		DialogDescription
	} from '$lib/components';
	import type { Organization } from '$lib/api/organization';
	import { API } from '$lib/api';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { queries } from '$lib/queries';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';
	import { goto } from '$app/navigation';

	interface DeleteOrgModal {
		org: Organization;
		open: boolean;
	}

	let { org, open = $bindable() }: DeleteOrgModal = $props();

	const queryClient = useQueryClient();

	const deleteOrgMutation = createMutation(() => ({
		mutationFn: () =>
			API.delete('/v1/orgs/{orgSlug}', {
				params: { orgSlug: org.slug }
			}),
		onSettled: () => {
			queryClient.invalidateQueries(queries.organization.list());
			queryClient.setQueryData(queries.organization.getByName(org.slug).queryKey, undefined);
			goto('/account/orgs');
		}
	}));
</script>

<Dialog bind:open>
	<DialogTitle>Delete <i>{org?.name}</i></DialogTitle>
	<DialogDescription>
		Are you sure you want to delete this organization? This action cannot be undone.
	</DialogDescription>
	{#if deleteOrgMutation.error}
		<Text variant="error">{deleteOrgMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			deleteOrgMutation.mutate(undefined, {
				onSuccess: () => {
					(e.target as HTMLFormElement)?.reset();
					open = false;
				}
			});
		}}
	>
		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button loading={deleteOrgMutation.isPending} color="red" class="self-end" type="submit">
				Delete Organization
			</Button>
		</DialogActions>
	</form>
</Dialog>
