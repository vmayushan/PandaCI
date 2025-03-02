<script lang="ts">
	import { page } from '$app/state';
	import { Code, Card, Heading } from '$lib/components';
	import Button from '$lib/components/button.svelte';
	import { authAPI } from '$lib/kratos';

	const errorId = page.url.searchParams.get('id');
	const errorObj = page.url.searchParams.get('error');

	let error: string | null = null;
	if (errorId && !errorObj) {
		authAPI
			.getFlowError({ id: errorId })
			.then((res) => {
				error = JSON.stringify(res.data.error, null, 2);
			})
			.catch((e) => {
				console.error(e);
				error = 'Failed to fetch error';
			});
	} else if (errorObj) {
		try {
			error = JSON.stringify(JSON.parse(decodeURIComponent(errorObj)), null, 2);
		} catch (e) {
			console.error(e);
			error = 'Failed to parse error';
		}
	} else {
		error = 'Unknown error';
	}
</script>

<Card class="mx-auto my-auto flex w-full max-w-md flex-col space-y-12 ">
	<Heading>Error</Heading>

	{#if error}
		<Code>{JSON.stringify(error)}</Code>
	{/if}

	<Button href="/logout?return_to=%2Flogin">Try again</Button>
</Card>
