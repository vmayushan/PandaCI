<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type LegendProps = SvelteHTMLElements['legend'];

	const { children, class: className, ...props }: LegendProps = $props();
</script>

<legend
	data-slot="legend"
	class={clsx(
		className,
		'data-disabled:opacity-50 text-base/6 font-semibold text-zinc-950 sm:text-sm/6 dark:text-white'
	)}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</legend>
