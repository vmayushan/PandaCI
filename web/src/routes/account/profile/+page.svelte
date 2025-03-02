<script lang="ts">
	import { page } from '$app/state';
	import { authAPI } from '$lib/kratos';
	import Messages from '$lib/components/kratos/messages.svelte';
	import Render from '$lib/components/kratos/render.svelte';
	import Divider from '$lib/components/divider.svelte';
	import {
		Button,
		DialogActions,
		DialogCloseButton,
		DialogDescription,
		DialogTitle,
		Fieldset,
		Text
	} from '$lib/components';
	import { handleError } from '../../(auth)/handleError.svelte';
	import { replaceState } from '$app/navigation';
	import Title from '$lib/components/title.svelte';
	import FieldGroup from '$lib/components/fieldset/fieldGroup.svelte';
	import { handleForm } from '$lib/utils';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { queries } from '$lib/queries';
	import type { UpdateSettingsFlowWithProfileMethod } from '@ory/client';
	import Field from '$lib/components/fieldset/field.svelte';
	import Label from '$lib/components/fieldset/label.svelte';
	import Skeleton from '$lib/components/skeleton.svelte';
	import Dialog from '$lib/components/dialog/dialog.svelte';
	import Card from '$lib/components/card.svelte';

	const flowId = page.url.searchParams.get('flow');

	const user = createQuery(() => ({
		...queries.auth.session()
	}));

	let reauthenticateModalOpen = $state(false);
	let reauthenticate = $state(false);

	const queryClient = useQueryClient();

	$effect(() => {
		if (
			user.data?.data?.authenticated_at &&
			new Date(user.data?.data.authenticated_at).getTime() < new Date().getTime() - 1000 * 60 * 14
		) {
			// Kratos requires a login 15 mins ago, if we leave 14 minutes then a user has 1 min to make changes
			reauthenticateModalOpen = reauthenticate ? reauthenticateModalOpen : true;
			reauthenticate = true;
		} else if (user.data?.data?.authenticated_at) {
			reauthenticateModalOpen = false;
			reauthenticate = false;

			setTimeout(
				() => {
					queryClient.refetchQueries(queries.auth.session());
				},
				new Date(user.data.data.authenticated_at).getTime() + 1000 * 60 * 14 - new Date().getTime()
			);
		}
	});

	let session = $state(
		flowId
			? authAPI
					.getSettingsFlow({
						id: flowId
					})
					.then((res) => {
						if (res.data.continue_with) {
							for (const continueWith of res.data.continue_with) {
								if (continueWith.action === 'show_verification_ui' && continueWith.flow.url) {
									window.location.href = continueWith.flow.url;
								}
							}
						}

						return res;
					})
					.catch(handleError)
			: authAPI
					.createBrowserSettingsFlow({})
					.then(async (res) => {
						page.url.searchParams.set('flow', res.data.id);
						replaceState(page.url, page.state);
						return res;
					})
					.catch(handleError)
	);
</script>

<Dialog bind:open={reauthenticateModalOpen}>
	<DialogTitle>Reauthenticate</DialogTitle>
	<DialogDescription>To make changes to your profile, you need to login again.</DialogDescription>

	<DialogActions>
		<DialogCloseButton plain>Cancel</DialogCloseButton>
		<Button autofocus color="dark/white" href="/login?refresh=true&return_to=/account/profile">
			Reauthenticate
		</Button>
	</DialogActions>
</Dialog>

<Title title="Profile">
	{#snippet description()}
		Update your profile settings
	{/snippet}
</Title>

{#if reauthenticate}
	<Card
		class="mt-4 flex flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0"
	>
		<Text variant="emphasis">To make changes to your profile, you need to login again.</Text>
		<Button class="mt-2" href="/login?refresh=true&return_to=/account/profile" color="dark/white">
			Reauthenticate
		</Button>
	</Card>
{/if}

<div class="mx-auto mt-12 max-w-3xl">
	{#await session}
		<Fieldset>
			<FieldGroup>
				{#each ['Email', 'Name', 'Avatar URL'] as field (field)}
					<Field>
						<Label>
							{field}
						</Label>
						<Skeleton class="my-2 h-10 w-full" />
					</Field>
				{/each}
				<Skeleton class="h-9 w-14" />
			</FieldGroup>
		</Fieldset>
	{:then sessionData}
		{#if sessionData}
			{@const data = sessionData.data}
			{#if data.ui.messages}
				<Messages messages={data.ui.messages} />
				<Divider class="mb-8 mt-4" />
			{/if}
			<form
				onsubmit={(e) => {
					const { data: body } = handleForm<UpdateSettingsFlowWithProfileMethod>(e);

					authAPI
						.updateSettingsFlow({
							flow: data.id,
							updateSettingsFlowBody: { ...body, method: 'profile' }
						})
						.then((res) => {
							session = Promise.resolve(res);
							queryClient.invalidateQueries(queries.auth.session());

							if (res.data.continue_with) {
								for (const continueWith of res.data.continue_with) {
									if (continueWith.action === 'show_verification_ui' && continueWith.flow.url) {
										window.location.href = continueWith.flow.url;
									}
								}
							}
						})
						.catch(handleError)
						.catch((err) => {
							if (err?.response) {
								session = Promise.resolve(err.response);
							}
						});
				}}
			>
				<Fieldset>
					<FieldGroup>
						<Render disabled={reauthenticate} smallButtons nodes={data.ui.nodes} />
					</FieldGroup>
				</Fieldset>
			</form>
		{/if}
	{/await}
</div>
