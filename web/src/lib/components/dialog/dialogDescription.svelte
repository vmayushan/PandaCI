<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { Text } from '../text';

	type DialogDescriptionProps = SvelteHTMLElements['p'];

	const { children, class: className, ...props }: DialogDescriptionProps = $props();
</script>

<Text class={[className, 'mt-2 text-pretty']} {...props}>
	{#if children}
		{@render children()}
	{/if}
</Text>
