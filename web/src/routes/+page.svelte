<script lang="ts">
	import { goto } from '$app/navigation';
	import { getCurrentOrg, getUser } from '$lib/runes/index.svelte';
	import { page } from '$app/state';

	// Makes sure we redirect to the login page if the user is not logged in
	getUser();

	const currentOrg = getCurrentOrg({
		page
	});

	$effect.pre(() => {
		if (currentOrg.data) {
			goto(currentOrg.data.slug, { replaceState: true });
		}
	});
</script>
