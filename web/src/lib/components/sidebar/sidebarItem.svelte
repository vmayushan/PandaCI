<script lang="ts" generics="T extends string | undefined = undefined">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import TouchTarget from '../touchTarget.svelte';
	import { getContext } from 'svelte';
	import type { SidebarSectionContext } from './sidebarSection.svelte';
	import type { SidebarContext } from './sidebar.svelte';
	import { Tooltip } from '../tooltip';

	type SidebarItemProps<T> = {
		current?: boolean;
		href?: T;
		tooltip?: string;
		innerClass?: string;
	} & (T extends string ? SvelteHTMLElements['a'] : SvelteHTMLElements['button']);

	const {
		class: className,
		innerClass,
		current,
		children,
		tooltip,
		...props
	}: SidebarItemProps<T> = $props();

	const context = getContext<SidebarContext>('sidebar');

	const classes = $derived(
		clsx(
			// Base
			'flex cursor-pointer relative items-center gap-3 rounded-lg px-2 py-2.5 text-left text-base/6 font-medium text-zinc-950 sm:py-2 sm:text-sm/5',
			// Leading icon/icon-only
			'*:data-[slot=icon]:size-6 *:data-[slot=icon]:shrink-0 *:data-[slot=icon]:text-zinc-500 sm:*:data-[slot=icon]:size-5',
			// Trailing icon (down chevron or similar)
			'*:last:data-[slot=icon]:ml-auto *:last:data-[slot=icon]:size-5 sm:*:last:data-[slot=icon]:size-4',
			// Avatar
			'*:data-[slot=avatar]:-m-0.5 *:data-[slot=avatar]:size-7 *:data-[slot=avatar]:[--ring-opacity:10%] sm:*:data-[slot=avatar]:size-6',
			// Hover
			'hover:bg-zinc-950/5 hover:*:data-[slot=icon]:text-zinc-950',
			// Focus
			'focus:outline-hidden focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500',
			// Active
			'active:bg-zinc-950/5 active:*:data-[slot=icon]:text-zinc-950',
			// Current
			'data-current:*:data-[slot=icon]:text-zinc-950',
			// Dark mode
			'dark:text-white dark:*:data-[slot=icon]:text-zinc-400',
			'dark:hover:bg-white/5 dark:hover:*:data-[slot=icon]:text-white',
			'dark:active:bg-white/5 dark:active:*:data-[slot=icon]:text-white',
			'dark:data-current:*:data-[slot=icon]:text-white',
			// Collapsed
			context.collapsed ? 'w-fit' : 'w-full'
		)
	);

	const { animationName } = getContext<SidebarSectionContext>('sidebar-section') ?? {};

	const element = props.href ? 'a' : 'button';
</script>

<Tooltip disable={!context.collapsed} text={tooltip} class={clsx(className, 'relative')}>
	{#if current}
		<span
			style:view-transition-name={animationName}
			class="indicator absolute inset-y-2 -left-4 w-0.5 rounded-full bg-zinc-950 dark:bg-white"
		></span>
	{/if}

	<svelte:element
		this={element}
		class={clsx(classes, innerClass)}
		data-current={current ? 'true' : undefined}
		{...props as any}
	>
		<TouchTarget />
		{@render children?.()}
	</svelte:element>
</Tooltip>
