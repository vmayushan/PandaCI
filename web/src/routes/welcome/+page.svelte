<script lang="ts">
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import CreateOrgForm, { type CreateOrgFormData } from './createOrgForm.svelte';
	import { Button, Card, Heading, Text, Divider, TextLink } from '$lib/components';
	import { X } from 'phosphor-svelte';
	import { goto } from '$app/navigation';
	import { useCreateOrgMutation } from '../account/useCreateOrg';
	import { handleForm } from '$lib/utils';

	const orgs = createQuery(() => queries.organization.list());

	$effect(() => {
		if (orgs.data && orgs.data.length) {
			goto('/');
		}
	});

	localStorage.setItem('welcome-visited', 'true');

	const createOrgMutation = useCreateOrgMutation();
</script>

<div class="flex h-full min-h-svh flex-col overflow-auto py-8">
	<Card class="relative m-auto w-full max-w-md">
		<div class="absolute right-4 top-4">
			<Button href="/" plain class=""><X class="h-4 w-4" /></Button>
		</div>
		<Heading class="mb-2">Create your organization</Heading>
		<Text>Set up your organization's workspace on our platform.</Text>

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
			<Button loading={createOrgMutation.isPending} color="dark/white" full type="submit">
				Create organization
			</Button>
		</form>

		<Divider class="mb-12 mt-8" />

		<Text class="text-center">
			Want to join an organization? <TextLink
				href="https://pandaci.com/docs/platform/other/org-members"
			>
				Ask your admin for an invite
			</TextLink>
		</Text>
	</Card>
</div>
