<script lang="ts">
	import { Heading, Button, Text } from '$lib/components';

	interface NotFoundProps {
		itemName: string;
	}
	const { itemName }: NotFoundProps = $props();
</script>

<main
	class="grid min-h-svh place-items-center bg-white px-6 py-24 sm:py-32 lg:px-8 dark:bg-zinc-900"
>
	<div class="text-center">
		<Text>404</Text>
		<Heading>{itemName} not found</Heading>
		<Text>Sorry, we couldn’t find the page you’re looking for.</Text>
		<div class="mt-10 flex items-center justify-center gap-x-6">
			<Button href="/">Go back home</Button>
		</div>
	</div>
</main>
