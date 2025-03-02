<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DescriptionTermProps = SvelteHTMLElements['dt'] & {};

	const { children, class: className, ...props }: DescriptionTermProps = $props();
</script>

<dt
	class={clsx(
		'col-start-1 border-t border-zinc-950/5 pt-3 text-zinc-500 first:border-none sm:border-t sm:border-zinc-950/5 sm:py-3 dark:border-white/5 dark:text-zinc-400 sm:dark:border-white/5',
		className
	)}
	{...props}
>
	{@render children?.()}
</dt>
