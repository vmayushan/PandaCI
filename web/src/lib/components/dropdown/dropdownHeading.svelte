<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DropdownHeadingProps = SvelteHTMLElements['header'];

	const { children, class: className, ...props }: DropdownHeadingProps = $props();
</script>

<header
	{...props}
	role="presentation"
	class={clsx(
		className,
		'col-span-full grid grid-cols-[1fr_auto] gap-x-12 px-3.5 pb-1 pt-2 text-sm/5 font-medium text-zinc-500 sm:px-3 sm:text-xs/5 dark:text-zinc-400'
	)}
>
	{@render children?.()}
</header>
