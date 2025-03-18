<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { API, type apiDefs, type Organization } from '$lib/api';
	import {
		Fieldset,
		FieldGroup,
		Field,
		Label,
		Input,
		Description,
		Button,
		Divider
	} from '$lib/components';
	import Text from '$lib/components/text/text.svelte';
	import Title from '$lib/components/title.svelte';
	import { queries } from '$lib/queries';
	import { handleForm } from '$lib/utils';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import DeleteOrgModal from './deleteOrgModal.svelte';

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	let orgURL = $state('');

	$effect(() => {
		orgURL = org.data?.slug || '';
	});

	const queryClient = useQueryClient();

	const updateOrgMutation = createMutation(() => ({
		mutationFn: (data: apiDefs['PUT']['/v1/orgs/{orgName}']['req']) =>
			API.put('/v1/orgs/{orgName}', {
				body: { ...data },
				params: { orgName: page.params.orgName }
			}),
		onSuccess: (updatedOrg) => {
			queryClient.setQueryData(
				queries.organization.getByName(updatedOrg.slug).queryKey,
				updatedOrg
			);
			queryClient.setQueryData(queries.organization.list().queryKey, (oldData?: Organization[]) => {
				if (!oldData) return [updatedOrg];

				const newOrgs = oldData.filter(
					(org) => org.slug !== updatedOrg.slug || org.slug === page.params.orgName
				);
				newOrgs.push(updatedOrg);
				return newOrgs;
			});
			queryClient.invalidateQueries(queries.organization.getByName(page.params.orgName));
			queryClient.invalidateQueries(queries.organization.list());
			if (updatedOrg.slug !== org.data?.slug) {
				goto(`/${updatedOrg.slug}/~/settings`);
			}
		}
	}));

	let deleteOrgModalOpen = $state(false);
</script>

{#if org.data}
	<DeleteOrgModal org={org.data} bind:open={deleteOrgModalOpen} />
{/if}

<Title title="Organization settings">
	{#snippet description()}
		Update your organization settings.
	{/snippet}
</Title>

{#if updateOrgMutation.isError}
	<Text class="mx-auto max-w-3xl py-2" variant="error">
		{updateOrgMutation.error.message}
	</Text>
{/if}

<form
	class="mx-auto mt-10 max-w-3xl"
	onsubmit={(e) => {
		const { data } = handleForm<apiDefs['PUT']['/v1/orgs/{orgName}']['req']>(e);

		updateOrgMutation.mutate(data);
	}}
>
	<Fieldset class="flex flex-col">
		<FieldGroup>
			<Field>
				<Label>Organization Name</Label>
				<Input defaultValue={org.data?.name} name="name" />
			</Field>

			<Field>
				<Label>Organization URL</Label>
				<Input
					oninput={(e) => {
						orgURL = e.currentTarget.value.replaceAll(' ', '-');
						e.currentTarget.value = orgURL;
					}}
					minlength={2}
					type="text"
					required
					placeholder="your-org"
					name="slug"
					defaultValue={org.data?.slug}
				/>
				<Description>
					https://app.pandaci.com/<b>
						{encodeURIComponent(orgURL || 'your-org')}
					</b>
				</Description>
			</Field>

			<Field>
				<Label>Organization Avatar URL</Label>
				<Input defaultValue={org.data?.avatarURL} name="avatarURL" />
			</Field>
		</FieldGroup>

		<Button
			loading={updateOrgMutation.isPending}
			type="submit"
			class="mt-8 self-end"
			color="dark/white">Save</Button
		>
	</Fieldset>
</form>

<Divider class="my-8" />

<div class="flex max-w-5xl">
	<Button
		onclick={() => {
			deleteOrgModalOpen = true;
		}}
		color="red"
		class="ml-auto"
	>
		Delete organization
	</Button>
</div>
