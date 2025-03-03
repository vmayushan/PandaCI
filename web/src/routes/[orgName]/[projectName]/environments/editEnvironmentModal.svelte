<script lang="ts">
	import {
		Dialog,
		DialogTitle,
		DialogBody,
		DialogActions,
		Text,
		Fieldset,
		Field,
		FieldGroup,
		Label,
		Input,
		Button,
		Description
	} from '$lib/components';
	import type { Project, Organization, ProjectEnvironment } from '$lib/api/organization';
	import { API, type apiDefs } from '$lib/api';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { handleForm } from '$lib/utils';
	import { queries } from '$lib/queries';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';
	import { untrack } from 'svelte';

	interface CreateVariableModal {
		project: Project;
		org: Organization;
		oldEnvironment?: ProjectEnvironment;
	}

	let { oldEnvironment = $bindable(), org, project }: CreateVariableModal = $props();

	const queryClient = useQueryClient();

	const projectEnvironmentMutation = createMutation(() => ({
		mutationFn: (
			data: apiDefs['PUT']['/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}']['req']
		) =>
			API.put('/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}', {
				body: { ...data },
				params: { orgSlug: org.slug, projectSlug: project.slug, environmentID: oldEnvironment!.id }
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.environments.projectEnvironments(org.slug, project.slug)._ctx.list()
			);
		}
	}));

	let open = $state(oldEnvironment !== undefined);

	$effect(() => {
		if (oldEnvironment) {
			untrack(() => {
				open = true;
			});
		}
	});

	$effect(() => {
		if (!open) {
			untrack(() => {
				oldEnvironment = undefined;
			});
		}
	});
</script>

<Dialog bind:open>
	<DialogTitle>Edit environment</DialogTitle>
	{#if projectEnvironmentMutation.error}
		<Text variant="error">{projectEnvironmentMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			const { data } =
				handleForm<
					apiDefs['PUT']['/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}']['req']
				>(e);

			projectEnvironmentMutation.mutate(
				{
					...data
				},
				{
					onSuccess: () => {
						(e.target as HTMLFormElement)?.reset();
						open = false;
					}
				}
			);
		}}
	>
		<DialogBody>
			<Fieldset>
				<FieldGroup>
					<Field>
						<Label>Name</Label>
						<Input defaultValue={oldEnvironment?.name} autofocus={true} type="text" name="name" />
					</Field>

					<Field>
						<Label>Branch Pattern</Label>
						<Description>
							Use globs to match branches, e.g <pre class="inline">feat/*</pre>
						</Description>
						<Input defaultValue={oldEnvironment?.branchPattern} type="text" name="branchPattern" />
					</Field>
				</FieldGroup>
			</Fieldset>
		</DialogBody>
		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button
				loading={projectEnvironmentMutation.isPending}
				color="dark/white"
				class="self-end"
				type="submit"
			>
				Save
			</Button>
		</DialogActions>
	</form>
</Dialog>
