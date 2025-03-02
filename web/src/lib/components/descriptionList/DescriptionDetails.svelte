<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DescriptionDetailsProps = SvelteHTMLElements['dd'] & {};

	const { children, class: className, ...props }: DescriptionDetailsProps = $props();
</script>

<dd
	class={clsx(
		'pb-3 pt-1 text-zinc-950 sm:border-t sm:border-zinc-950/5 sm:py-3 dark:text-white dark:sm:border-white/5 sm:[&:nth-child(2)]:border-none',
		className
	)}
	{...props}
>
	{@render children?.()}
</dd>
