<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { API, type apiDefs } from '$lib/api';
	import type { Project } from '$lib/api/organization';
	import {
		Input,
		Label,
		Text,
		Fieldset,
		FieldGroup,
		Description,
		Title,
		Button
	} from '$lib/components';
	import Divider from '$lib/components/divider.svelte';
	import Field from '$lib/components/fieldset/field.svelte';
	import { queries } from '$lib/queries';
	import { handleForm } from '$lib/utils';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';

	let projectURL = $state('');

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	$effect(() => {
		projectURL = project.data?.slug || '';
	});

	const queryClient = useQueryClient();

	const updateProjectMutation = createMutation(() => ({
		mutationFn: (data: apiDefs['PUT']['/v1/orgs/{orgName}/projects/{projectName}']['req']) =>
			API.put('/v1/orgs/{orgName}/projects/{projectName}', {
				body: { ...data },
				params: { orgName: page.params.orgName, projectName: page.params.projectName }
			}),
		onMutate: (data) => {
			queryClient.setQueryData(
				queries.organization
					.getByName(page.params.orgName)
					._ctx.projectByName(page.params.projectName).queryKey,
				(old?: Project) => ({ ...old, ...data })
			);

			queryClient.setQueryData(
				queries.organization.getByName(page.params.orgName)._ctx.projects().queryKey,
				(old?: Project[]) =>
					old?.map((project) =>
						project.slug === page.params.projectName ? { ...project, ...data } : project
					)
			);
		},
		onSuccess: (updatedProject) => {
			queryClient.setQueryData(
				queries.organization.getByName(page.params.orgName)._ctx.projectByName(updatedProject.slug)
					.queryKey,
				updatedProject
			);
			if (updatedProject.slug !== project.data?.slug) {
				goto(`/${page.params.orgName}/${updatedProject.slug}/settings`);
			}
		},
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.organization.getByName(page.params.orgName)._ctx.projects()
			);
		}
	}));

	const deleteProjectMutation = createMutation(() => ({
		mutationFn: () =>
			API.delete('/v1/orgs/{orgName}/projects/{projectName}', {
				params: { orgName: page.params.orgName, projectName: page.params.projectName }
			}),
		onSuccess: () => {
			queryClient.invalidateQueries(queries.organization.getByName(page.params.orgName));
			goto(`/${page.params.orgName}`);
		}
	}));
</script>

<Title title="Project settings">
	{#snippet description()}
		Update your project settings.
	{/snippet}
</Title>

{#if updateProjectMutation.error}
	<Text variant="error">
		{updateProjectMutation.error.message}
	</Text>
{/if}

<form
	class="mx-auto mt-10 max-w-3xl"
	onsubmit={(e) => {
		const { data } =
			handleForm<apiDefs['PUT']['/v1/orgs/{orgName}/projects/{projectName}']['req']>(e);

		updateProjectMutation.mutate(data);
	}}
>
	<Fieldset class="flex flex-col">
		<FieldGroup>
			<Field>
				<Label>Project Name</Label>
				<Input defaultValue={project.data?.name} name="name" />
			</Field>

			<Field>
				<Label>Project URL</Label>
				<Input
					oninput={(e) => {
						projectURL = e.currentTarget.value.replaceAll(' ', '-');
						e.currentTarget.value = projectURL;
					}}
					minlength={2}
					type="text"
					required
					placeholder="your-org"
					name="slug"
					defaultValue={project.data?.slug}
				/>
				<Description>
					https://app.pandaci.com/{page.params.orgName}/<b>
						{encodeURIComponent(projectURL || 'your-project')}
					</b>
				</Description>
			</Field>

			<Field>
				<Label>Project Avatar URL</Label>
				<Input defaultValue={project.data?.avatarURL} name="avatarURL" />
			</Field>
		</FieldGroup>

		<Button
			loading={updateProjectMutation.isPending}
			type="submit"
			class="mt-8 self-end"
			color="dark/white">Save</Button
		>
	</Fieldset>
</form>

<Divider class="my-8" />

<div class="flex max-w-5xl">
	<Button
		loading={deleteProjectMutation.isPending}
		onclick={() => deleteProjectMutation.mutate()}
		color="red"
		class="ml-auto"
	>
		Delete project
	</Button>
</div>
