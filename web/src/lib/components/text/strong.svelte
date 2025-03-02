<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type StrongProps = SvelteHTMLElements['strong'];

	const { children, class: className, ...props }: StrongProps = $props();
</script>

<strong {...props} class={clsx(className, 'font-medium text-zinc-950 dark:text-white')}>
	{#if children}
		{@render children()}
	{/if}
</strong>
