<script lang="ts">
	import { portal } from '@zag-js/svelte';
	import { getContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { DropdownContext } from './dropdown.svelte';

	type DropdownItemProps = SvelteHTMLElements['div'] & {
		disablePortal?: boolean;
	};

	const { children, class: className, disablePortal, ...props }: DropdownItemProps = $props();

	const subCtx = getContext<DropdownContext>('dropdown-sub');
</script>

<div use:portal={{ disabled: disablePortal }} {...subCtx.api.getPositionerProps()}>
	<div
		{...props}
		{...subCtx.api.getContentProps()}
		class={[
			className,
			// Base styles
			'isolate z-10 w-max min-w-full rounded-xl p-1',
			// Invisible border that is only visible in `forced-colors` mode for accessibility purposes
			'focus:outline-hidden outline outline-transparent',
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
