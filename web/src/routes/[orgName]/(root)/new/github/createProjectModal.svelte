<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { API, type apiDefs, type GitInstallation, type GitRepository } from '$lib/api';
	import {
		Button,
		Description,
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
	import Text from '$lib/components/text/text.svelte';
	import { queries } from '$lib/queries';
	import { handleForm } from '$lib/utils';
	import { QueryClient, createMutation } from '@tanstack/svelte-query';

	interface CreateProjectModal {
		open: boolean;
		repo: GitRepository;
		installation: GitInstallation;
	}

	let { open = $bindable(), repo, installation }: CreateProjectModal = $props();

	type ProjectFormData = apiDefs['POST']['/v1/orgs/{orgName}/projects']['req'];

	let projectURL = $state<ProjectFormData['slug']>();
	let displayName = $state<ProjectFormData['name']>();

	let displayNameDirty = $state(false);

	$effect(() => {
		if (!displayNameDirty) displayName = projectURL?.replaceAll('-', ' ');
	});

	const queryClient = new QueryClient();

	const createProjectMutation = createMutation(() => ({
		mutationFn: (data: apiDefs['POST']['/v1/orgs/{orgName}/projects']['req']) =>
			API.post('/v1/orgs/{orgName}/projects', {
				body: { ...data },
				params: { orgName: page.params.orgName }
			}),
		onSuccess: (project) => {
			queryClient.setQueryData(
				queries.organization.getByName(page.params.orgName)._ctx.projectByName(project.slug)
					.queryKey,
				project
			);
			queryClient.setQueryData(
				queries.runs.projectRuns(page.params.orgName, project.slug)._ctx.list().queryKey,
				[]
			);

			goto(`/${page.params.orgName}/${project.slug}`);
		}
	}));

	$effect(() => {
		if (repo) {
			projectURL = repo.name.replaceAll(' ', '-');
			displayName = repo.name;
			displayNameDirty = false;
		}
	});

	$effect(() => {
		if (!open) {
			projectURL = repo.name.replaceAll(' ', '-');
			displayName = repo.name;
			displayNameDirty = false;
		}
	});
</script>

<Dialog bind:open>
	<DialogTitle>Create a project</DialogTitle>
	<DialogDescription>Set up a new CI/CD pipeline for your application.</DialogDescription>

	{#if createProjectMutation.error}
		<Text variant="error">{createProjectMutation.error.message}</Text>
	{/if}

	<form
		onsubmit={(e) => {
			const { data } =
				handleForm<Pick<apiDefs['POST']['/v1/orgs/{orgName}/projects']['req'], 'slug' | 'name'>>(e);

			createProjectMutation.mutate({
				...data,
				gitProviderType: installation.type,
				gitProviderIntegrationID: installation.id,
				gitProviderRepoID: repo.id
			});
		}}
	>
		<DialogBody>
			<Fieldset>
				<FieldGroup>
					<Field>
						<Label>Project URL</Label>
						<Input
							defaultValue={projectURL}
							id="create-project-url"
							oninput={(e) => {
								projectURL = e.currentTarget.value.replaceAll(' ', '-');
								e.currentTarget.value = projectURL;
							}}
							minlength={2}
							type="text"
							required
							placeholder="your-org"
							name="slug"
						/>
						<Description>
							https://app.pandaci.com/{page.params.orgName}/<b>
								{encodeURIComponent(projectURL || 'your-project')}
							</b>
						</Description>
					</Field>

					<Field>
						<Label>Project Name</Label>
						<Input
							onchange={() => {
								displayNameDirty = displayName !== '';
								displayName = displayName === '' ? projectURL : displayName;
							}}
							required
							value={displayName}
							oninput={(e) => (displayName = e.currentTarget.value)}
							placeholder="your project"
							name="name"
							type="text"
						/>
					</Field>
				</FieldGroup>
			</Fieldset>
		</DialogBody>

		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button loading={createProjectMutation.isPending} color="dark/white" type="submit">
				Create project
			</Button>
		</DialogActions>
	</form>
</Dialog>
