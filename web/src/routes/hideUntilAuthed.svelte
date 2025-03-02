<script lang="ts">
	import { getUser } from '$lib/runes/index.svelte';
	import { page } from '$app/state';

	const user = getUser({ enabled: !page.route.id?.startsWith('/(auth)') });

	$effect.pre(() => {
		if (user.isLoading && !page.route.id?.startsWith('/(auth)')) {
			document.body.classList.add('opacity-0');
		} else {
			document.body.classList.remove('opacity-0');
		}
	});
</script>
