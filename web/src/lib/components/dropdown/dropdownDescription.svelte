<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DropdownItemProps = SvelteHTMLElements['p'];

	const { children, class: className, ...props }: DropdownItemProps = $props();
</script>

<p
	data-slot="description"
	class={clsx(
		'group-data-highlighted:text-white forced-colors:group-data-highlighted:text-[HighlightText] col-span-2 col-start-2 row-start-2 text-sm/5 text-zinc-500 sm:text-xs/5 dark:text-zinc-400',
		className
	)}
>
	{@render children?.()}
</p>
