<script lang="ts">
	import { page } from '$app/state';
	import { authAPI } from '$lib/kratos';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { handleError } from '../handleError.svelte';
	import { queries } from '$lib/queries';

	const returnTo = page.url.searchParams.get('return_to') || undefined;

	const queryClient = useQueryClient();

	authAPI
		.createBrowserLogoutFlow({
			returnTo
		})
		.then((res) => {
			queryClient.invalidateQueries(queries.auth.session());
			window.location.href = res.data.logout_url;
		})
		.catch((err) => {
			if (err.response.status === 404) {
				window.location.href = '/';
			}
			handleError(err);
		});
</script>
