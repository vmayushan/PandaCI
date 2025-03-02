<script lang="ts">
	import clsx from 'clsx';
	import { getContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { SidebarContext } from './sidebar.svelte';

	type SidebarHeadingProps = SvelteHTMLElements['h3'];

	const { class: className, children, ...props }: SidebarHeadingProps = $props();

	const context = getContext<SidebarContext>('sidebar');
</script>

<h3
	class={clsx(
		'mb-1 px-2 text-xs/6 font-medium text-zinc-500 transition dark:text-zinc-400',
		context.collapsed && 'opacity-0',
		className
	)}
	{...props}
>
	{@render children?.()}
</h3>
