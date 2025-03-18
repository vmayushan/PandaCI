<script lang="ts">
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import CreateOrgForm, { type CreateOrgFormData } from './createOrgForm.svelte';
	import { Button, Card, Heading, Text, Divider, TextLink } from '$lib/components';
	import { goto } from '$app/navigation';
	import { useCreateOrgMutation } from '../account/useCreateOrg';
	import { handleForm } from '$lib/utils';

	const orgs = createQuery(() => queries.organization.list());

	let loaded = $state(false);

	$effect(() => {
		if (!loaded && orgs.data && orgs.data.length) {
			// goto('/');
		} else if (orgs.isFetched) {
			loaded = true;
		}
	});

	localStorage.setItem('welcome-visited', 'true');

	const createOrgMutation = useCreateOrgMutation();
</script>

<div class="sm:bg-surface-low flex h-full min-h-svh flex-col overflow-auto py-8">
	<Card hideMobileRing class="relative m-auto w-full max-w-md !ring-0 sm:ring-1">
		<Heading class="mb-2">Create your organization</Heading>
		<Text>Set up your organization's (or personal) workspace on our platform.</Text>

		{#if createOrgMutation.error}
			<Text variant="error">{createOrgMutation.error.message}</Text>
		{/if}
		<form
			onsubmit={(e) => {
				const { data } = handleForm<CreateOrgFormData>(e);

				createOrgMutation.mutate(data);
			}}
			class="mt-12 space-y-12"
		>
			<CreateOrgForm />
			<div class="flex flex-col items-center justify-end gap-3 sm:flex-row">
				<Button loading={createOrgMutation.isPending} color="dark/white" full type="submit">
					Create organization
				</Button>

				<Button class="w-full sm:w-auto" plain href="/account">Skip</Button>
			</div>
		</form>

		<Divider class="my-8" />
		<Text>
			Want to join an org? <TextLink href="https://pandaci.com/docs/platform/other/org-members">
				Get an invite
			</TextLink>
		</Text>
	</Card>
</div>
