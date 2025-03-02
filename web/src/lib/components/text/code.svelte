<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type TextProps = SvelteHTMLElements['code'];

	const { children, class: className, ...props }: TextProps = $props();
</script>

<code
	{...props}
	class={clsx(
		className,
		'rounded-sm border border-zinc-950/10 bg-zinc-950/[2.5%] px-0.5 text-sm font-medium text-zinc-950 sm:text-[0.8125rem] dark:border-white/20 dark:bg-white/5 dark:text-white'
	)}
>
	{#if children}
		{@render children()}
	{/if}
</code>
