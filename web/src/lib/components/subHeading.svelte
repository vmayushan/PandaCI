<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SubHeadingProps = SvelteHTMLElements['h2'] & { level?: 1 | 2 | 3 | 4 | 5 | 6 };

	const { level = 2, children, class: className, ...props }: SubHeadingProps = $props();
</script>

<svelte:element
	this={`h${level}`}
	class={clsx(className, 'text-base/7 font-semibold text-zinc-950 sm:text-sm/6 dark:text-white')}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</svelte:element>
