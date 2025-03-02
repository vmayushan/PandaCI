<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type ErrorProps = SvelteHTMLElements['p'];

	const { children, class: className, ...props }: ErrorProps = $props();
</script>

<p
	data-slot="error"
	class={clsx(
		className,
		'data-disabled:opacity-50 text-base/6 text-red-600 sm:text-sm/6 dark:text-red-500'
	)}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</p>
