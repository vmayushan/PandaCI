<script lang="ts">
	import '../app.css';
	import '@fontsource-variable/inter';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { browser } from '$app/environment';
	import { onNavigate } from '$app/navigation';
	import type { APIError } from '$lib/api';
	import ThemeProvider from '$lib/components/themeProvider.svelte';
	import posthog from 'posthog-js';
	import { PUBLIC_STAGE } from '$env/static/public';

	const posthogToken =
		PUBLIC_STAGE === 'prod'
			? 'phc_DR3pD0efwAHFCZJD0qb4vvKb1aMMupRrufHAOcGO1rX'
			: 'phc_p45o9ThJKsfyDOXM79eh2j1BpMkmvfJr2ofzs2bp1NM';

	$effect.root(() => {
		if (browser) {
			posthog.init(posthogToken, {
				api_host:
					PUBLIC_STAGE === 'prod' ? 'https://hedgehog.pandaci.com' : 'https://us.i.posthog.com',
				person_profiles: 'identified_only',
				autocapture: false
			});
		}
		return;
	});

	const { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser,
				retry: (_, err) => (err as APIError).status !== 404
			}
		}
	});

	onNavigate(() => {
		if (!document.startViewTransition) return;

		return new Promise((fulfill) => {
			document.startViewTransition(() => new Promise(fulfill as any));
		});
	});
</script>

<svelte:head>
	<title>PandaCI</title>
	<meta
		name="description"
		content="PandaCI is an open-source CI/CD platform where workflows are coded, not defined"
	/>
</svelte:head>

<QueryClientProvider client={queryClient}>
	<ThemeProvider>
		{@render children()}
	</ThemeProvider>
</QueryClientProvider>
