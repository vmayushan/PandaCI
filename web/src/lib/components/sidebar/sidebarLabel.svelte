<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { SidebarContext } from './sidebar.svelte';
	import { getContext } from 'svelte';

	type SidebarLabelProps = SvelteHTMLElements['span'];

	const { class: className, children, ...props }: SidebarLabelProps = $props();

	const context = getContext<SidebarContext>('sidebar');
</script>

<span class={clsx('truncate', className, context.collapsed && 'sr-only')} {...props}>
	{#if children}
		{@render children()}
	{/if}
</span>
