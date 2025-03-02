<script lang="ts">
	import { getContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { DropdownContext } from './dropdown.svelte';

	type DropdownItemProps = SvelteHTMLElements['span'];

	const { children, class: className, ...props }: DropdownItemProps = $props();

	const ctx = getContext<DropdownContext>('dropdown');

	const subCtx = getContext<DropdownContext>('dropdown-sub');

	const triggerProps = $derived(ctx.api.getTriggerItemProps(subCtx.api));
</script>

<span
	{...triggerProps}
	{...props}
	class={[
		className,
		// Base styles
		'data-highlighted:outline-hidden group w-full cursor-default rounded-lg px-3.5 py-2.5 sm:px-3 sm:py-1.5',
		// Text styles
		'text-left text-base/6 text-zinc-950 sm:text-sm/6 dark:text-white forced-colors:text-[CanvasText]',
		// Focus
		'data-highlighted:bg-blue-500 data-highlighted:text-white',
		// Disabled state
		'disabled:opacity-50',
		// Forced colors mode
		'forced-colors:data-highlighted:bg-[Highlight] forced-colors:data-highlighted:text-[HighlightText] forced-colors:data-highlighted:*:data-[slot=icon]:text-[HighlightText] forced-color-adjust-none',
		// Use subgrid when available but fallback to an explicit grid layout if not
		'col-span-full grid grid-cols-[auto_1fr_1.5rem_0.5rem_auto] items-center supports-[grid-template-columns:subgrid]:grid-cols-subgrid',
		// Icons
		'*:data-[slot=icon]:col-start-1 *:data-[slot=icon]:row-start-1 *:data-[slot=icon]:-ml-0.5 *:data-[slot=icon]:mr-2.5 *:data-[slot=icon]:size-5 sm:*:data-[slot=icon]:mr-2 sm:*:data-[slot=icon]:size-4',
		'data-highlighted:*:data-[slot=icon]:text-white dark:data-highlighted:*:data-[slot=icon]:text-white *:data-[slot=icon]:text-zinc-500 dark:*:data-[slot=icon]:text-zinc-400',
		// Avatar
		'*:data-[slot=avatar]:-ml-1 *:data-[slot=avatar]:mr-2.5 *:data-[slot=avatar]:size-6 sm:*:data-[slot=avatar]:mr-2 sm:*:data-[slot=avatar]:size-5'
	]}
>
	{#if children}
		{@render children()}
	{/if}
</span>
