<script lang="ts">
	import { getContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { DropdownContext } from './dropdown.svelte';
	import type { OptionItemProps } from '@zag-js/menu';
	import CheckboxIndicator from '../checkbox/checkboxIndicator.svelte';

	type BaseDropdownItemProps = { href?: string; value: string; checked: boolean };

	type DropdownItemProps<T extends { href?: string } = object> = T extends { href: string }
		? BaseDropdownItemProps & SvelteHTMLElements['a']
		: BaseDropdownItemProps & SvelteHTMLElements['button'];

	let {
		children,
		href,
		class: className,
		value,
		checked = $bindable(),
		...props
	}: DropdownItemProps = $props();

	const ctx = getContext<DropdownContext>('dropdown');

	const item: OptionItemProps = $derived({
		value,
		type: 'checkbox',
		checked,
		closeOnSelect: false,
		onCheckedChange: (c) => (checked = c)
	});
</script>

<svelte:element
	this={href ? 'a' : 'button'}
	{href}
	type="button"
	tabindex="-1"
	{...ctx.api.getOptionItemProps(item)}
	{...props}
	class={[
		className,
		// Base styles
		'data-highlighted:outline-hidden group w-full cursor-default truncate rounded-lg px-3.5 py-2.5 sm:px-3 sm:py-1.5',
		// Text styles
		'text-left text-base/6 text-zinc-950 sm:text-sm/6 dark:text-white forced-colors:text-[CanvasText]',
		// Focus
		'focus:outline-hidden data-highlighted:bg-blue-500 data-highlighted:text-white focus:ring-2 focus:ring-blue-500',
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
	<span class="flex items-center space-x-2">
		<CheckboxIndicator />
		{#if children}
			<span>
				{@render children()}
			</span>
		{/if}
	</span>
</svelte:element>
