<script lang="ts">
	import { goto } from '$app/navigation';
	import { API, type apiDefs } from '$lib/api';
	import type { Organization, Project } from '$lib/api/organization';
	import {
		Button,
		Dialog,
		DialogActions,
		DialogBody,
		DialogCloseButton,
		DialogDescription,
		DialogTitle,
		Field,
		FieldGroup,
		Fieldset,
		Input,
		Label
	} from '$lib/components';
	import { handleForm } from '$lib/utils';
	import { createMutation } from '@tanstack/svelte-query';

	interface CreateProjectModal {
		open: boolean;
		project: Project;
		org: Organization;
	}

	let { open = $bindable(), org, project }: CreateProjectModal = $props();

	const triggerProjectMutation = createMutation(() => ({
		mutationFn: (
			data: apiDefs['POST']['/v1/orgs/{orgName}/projects/{projectName}/trigger']['req']
		) =>
			API.post('/v1/orgs/{orgName}/projects/{projectName}/trigger', {
				body: { ...data },
				params: { orgName: org.slug, projectName: project.slug }
			})
	}));
</script>

<Dialog bind:open>
	<DialogTitle>Trigger a run</DialogTitle>
	<DialogDescription>Runs are also triggered by git events such as a push</DialogDescription>
	<form
		onsubmit={(e) => {
			const { data } =
				handleForm<apiDefs['POST']['/v1/orgs/{orgName}/projects/{projectName}/trigger']['req']>(e);

			triggerProjectMutation.mutate(
				{
					...data
				},
				{
					onSuccess(data) {
						open = false;
						goto(`/${org.slug}/${project.slug}/runs/${data[0].number}`);
					}
				}
			);
		}}
	>
		<DialogBody>
			<Fieldset>
				<FieldGroup>
					<Field>
						<Label>Branch</Label>
						<Input type="text" required name="branch" />
					</Field>
					<Field>
						<Label>Commit</Label>
						<Input type="text" required name="sha" />
					</Field>
				</FieldGroup>
			</Fieldset>
		</DialogBody>

		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button loading={triggerProjectMutation.isPending} color="dark/white" type="submit">
				Trigger run
			</Button>
		</DialogActions>
	</form>
</Dialog>
