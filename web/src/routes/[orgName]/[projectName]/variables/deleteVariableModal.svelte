<script lang="ts">
	import {
		Dialog,
		DialogTitle,
		DialogActions,
		Text,
		Button,
		DialogDescription
	} from '$lib/components';
	import type { Project, Organization, ProjectVariable } from '$lib/api/organization';
	import { API } from '$lib/api';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { queries } from '$lib/queries';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';
	import { untrack } from 'svelte';

	interface CreateVariableModal {
		project: Project;
		org: Organization;
		variable?: ProjectVariable;
	}

	let { org, project, variable = $bindable() }: CreateVariableModal = $props();

	const queryClient = useQueryClient();

	const projectVariableMutation = createMutation(() => ({
		mutationFn: () =>
			API.delete('/v1/orgs/{orgSlug}/projects/{projectSlug}/variables/{variableID}', {
				params: { orgSlug: org.slug, projectSlug: project.slug, variableID: variable!.id }
			}),
		onSettled: () => {
			queryClient.setQueryData(
				queries.variables.projectVariables(org.slug, project.slug)._ctx.get(variable!.id).queryKey,
				undefined
			);
			queryClient.invalidateQueries(
				queries.variables.projectVariables(org.slug, project.slug)._ctx.list()
			);
			queryClient.invalidateQueries(
				queries.variables.projectVariables(org.slug, project.slug)._ctx.get(variable!.id)
			);
		}
	}));

	let open = $state(variable !== undefined);

	$effect(() => {
		if (variable) {
			untrack(() => {
				open = true;
			});
		}
	});

	$effect(() => {
		if (!open) {
			untrack(() => {
				variable = undefined;
			});
		}
	});
</script>

<Dialog bind:open>
	<DialogTitle>Delete <i>{variable?.key}</i></DialogTitle>
	<DialogDescription>
		Are you sure you want to delete this variable? This action cannot be undone.
	</DialogDescription>
	{#if projectVariableMutation.error}
		<Text variant="error">{projectVariableMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			projectVariableMutation.mutate(undefined, {
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
				loading={projectVariableMutation.isPending}
				color="red"
				class="self-end"
				type="submit"
			>
				Delete
			</Button>
		</DialogActions>
	</form>
</Dialog>
