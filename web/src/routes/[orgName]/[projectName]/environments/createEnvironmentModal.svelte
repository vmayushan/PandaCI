<script lang="ts">
	import {
		Dialog,
		DialogTitle,
		DialogDescription,
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
	import type { Project, Organization } from '$lib/api/organization';
	import { API, type apiDefs } from '$lib/api';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { handleForm } from '$lib/utils';
	import { queries } from '$lib/queries';
	import DialogCloseButton from '$lib/components/dialog/dialogCloseButton.svelte';

	interface CreateEnvironmentModal {
		open: boolean;
		project: Project;
		org: Organization;
	}

	let { open = $bindable(), org, project }: CreateEnvironmentModal = $props();

	const queryClient = useQueryClient();

	const projectEnvironmentMutation = createMutation(() => ({
		mutationFn: (
			data: apiDefs['POST']['/v1/orgs/{orgSlug}/projects/{projectSlug}/environments']['req']
		) =>
			API.post('/v1/orgs/{orgSlug}/projects/{projectSlug}/environments', {
				body: { ...data },
				params: { orgSlug: org.slug, projectSlug: project.slug }
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.environments.projectEnvironments(org.slug, project.slug)._ctx.list()
			);
		}
	}));
</script>

<Dialog bind:open>
	<DialogTitle>Create environment</DialogTitle>
	<DialogDescription>
		Environments are used to group variables and restrict which branches can access them.
	</DialogDescription>
	{#if projectEnvironmentMutation.error}
		<Text variant="error">{projectEnvironmentMutation.error.message}</Text>
	{/if}
	<form
		class="flex flex-col space-y-8"
		onsubmit={(e) => {
			const { data } =
				handleForm<
					apiDefs['POST']['/v1/orgs/{orgSlug}/projects/{projectSlug}/environments']['req']
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
						<Input autofocus={true} type="text" name="name" />
					</Field>

					<Field>
						<Label>Branch Pattern</Label>
						<Description>
							Use globs to match branches, e.g <pre class="inline">feat/*</pre>
						</Description>
						<Input type="text" name="branchPattern" />
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
