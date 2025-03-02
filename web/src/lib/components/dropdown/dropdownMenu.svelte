<script lang="ts">
	import { portal } from '@zag-js/svelte';
	import { getContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { DropdownContext } from './dropdown.svelte';

	type DropdownMenuProps = SvelteHTMLElements['div'] & {
		disablePortal?: boolean;
	};

	const { children, class: className, disablePortal, ...props }: DropdownMenuProps = $props();

	const ctx = getContext<DropdownContext>('dropdown');
</script>

<div
	class={[disablePortal && '[--z-index:100]!']}
	use:portal={{ disabled: disablePortal }}
	{...ctx.api.getPositionerProps()}
>
	<div
		{...props}
		{...ctx.api.getContentProps()}
		class={[
			className,
			// Base styles
			'isolate z-50 w-max min-w-full rounded-xl p-1',
			// Invisible border that is only visible in `forced-colors` mode for accessibility purposes
			'focus:outline-hidden outline-1 outline-transparent',
			// Handle scrolling when menu won't fit in viewport
			'overflow-y-auto',
			// Popover background
			'bg-white/75 backdrop-blur-xl dark:bg-zinc-800/75',
			// Shadows
			'shadow-lg ring-1 ring-zinc-950/10 dark:ring-inset dark:ring-white/10'
		]}
	>
		{#if children}
			{@render children()}
		{/if}
	</div>
</div>
