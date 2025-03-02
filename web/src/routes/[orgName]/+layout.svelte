<script lang="ts">
	import { page } from '$app/state';
	import type { APIError } from '$lib/api';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import type { Snippet } from 'svelte';
	import NotFound from './notFound.svelte';
	import PaddleContext from '$lib/components/paddle/paddleContext.svelte';
	import posthog from 'posthog-js';
	import { getUser } from '$lib/runes/user.svelte';

	const { children }: { children: Snippet } = $props();

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));
	const user = getUser();

	$effect(() => {
		if (user.data)
			posthog.identify(user.data.id, {
				email: user.data.email,
				name: user.data.name
			});
	});
</script>

<PaddleContext>
	{#if org.isError && (org.error as APIError).status === 404}
		<NotFound itemName="Org" />
	{:else}
		{@render children()}
	{/if}
</PaddleContext>
