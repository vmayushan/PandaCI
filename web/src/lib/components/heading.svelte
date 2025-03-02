<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type HeadingProps = SvelteHTMLElements['h1'] & {
		level?: 1 | 2 | 3 | 4 | 5 | 6;

		size?: 'lg' | 'base' | 'sm' | 'xs';
	};

	const { level = 1, children, class: className, size = 'base', ...props }: HeadingProps = $props();
</script>

<svelte:element
	this={`h${level}`}
	class={clsx(
		className,
		'font-semibold text-zinc-950 dark:text-white',
		size === 'lg' && 'text-3xl/8 sm:text-2xl/8',
		size === 'base' && 'text-2xl/8 sm:text-xl/8',
		size === 'sm' && 'text-xl/8 sm:text-lg/8'
	)}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</svelte:element>
