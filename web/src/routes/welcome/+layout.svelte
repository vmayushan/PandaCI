<script lang="ts">
	import { getUser } from '$lib/runes/user.svelte';
	import posthog from 'posthog-js';

	const { children } = $props();

	const user = getUser();

	$effect(() => {
		if (user.data) {
			posthog.identify(user.data.id, {
				email: user.data.email,
				name: user.data.name
			});
			posthog.resetGroups();
		}
	});
</script>

{@render children()}
