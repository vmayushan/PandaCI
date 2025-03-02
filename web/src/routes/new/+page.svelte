<script lang="ts">
	import { goto } from '$app/navigation';
	import { getCurrentOrg } from '$lib/runes/index.svelte';
	import { page } from '$app/state';

	const currentOrg = getCurrentOrg({
		page
	});

	$effect(() => {
		if (currentOrg.data) {
			goto(`/${currentOrg.data.slug}/new/github`, { replaceState: true });
		}
	});
</script>
