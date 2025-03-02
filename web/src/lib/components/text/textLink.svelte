<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type TextProps = SvelteHTMLElements['a'];

	const { children, class: className, ...props }: TextProps = $props();
</script>

<a
	{...props}
	class={clsx(
		className,
		'data-hover:decoration-zinc-950 dark:data-hover:decoration-white text-zinc-950 underline decoration-zinc-950/50 dark:text-white dark:decoration-white/50'
	)}
>
	{#if children}
		{@render children()}
	{/if}
</a>
