<script lang="ts">
	import {
		Dialog,
		DialogTitle,
		DialogActions,
		Text,
		Button,
		DialogDescription
	} from '$lib/components';
	import type { Project, Organization, ProjectEnvironment } from '$lib/api/organization';
	import { API } from '$lib/api';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { queries } from '$lib/queries';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';
	import { untrack } from 'svelte';

	interface DeleteEnvironmentModal {
		project: Project;
		org: Organization;
		environment?: ProjectEnvironment;
	}

	let { org, project, environment = $bindable() }: DeleteEnvironmentModal = $props();

	const queryClient = useQueryClient();

	const projectEnvironmentMutation = createMutation(() => ({
		mutationFn: () =>
			API.delete('/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}', {
				params: { orgSlug: org.slug, projectSlug: project.slug, environmentID: environment!.id }
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.environments.projectEnvironments(org.slug, project.slug)._ctx.list()
			);
			queryClient.invalidateQueries(
				queries.variables.projectVariables(org.slug, project.slug)._ctx.list()
			);
		}
	}));

	let open = $state(environment !== undefined);

	$effect(() => {
		if (environment) {
			untrack(() => {
				open = true;
			});
		}
	});

	$effect(() => {
		if (!open) {
			untrack(() => {
				environment = undefined;
			});
		}
	});
</script>

<Dialog bind:open>
	<DialogTitle>Delete <i>{environment?.name}</i></DialogTitle>
	<DialogDescription>
		Are you sure you want to delete this environment? This action cannot be undone.
	</DialogDescription>
	{#if projectEnvironmentMutation.error}
		<Text variant="error">{projectEnvironmentMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			projectEnvironmentMutation.mutate(undefined, {
				onSuccess: () => {
					(e.target as HTMLFormElement)?.reset();
					open = false;
				}
			});
		}}
	>
		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button
				loading={projectEnvironmentMutation.isPending}
				color="red"
				class="self-end"
				type="submit"
			>
				Delete
			</Button>
		</DialogActions>
	</form>
</Dialog>
