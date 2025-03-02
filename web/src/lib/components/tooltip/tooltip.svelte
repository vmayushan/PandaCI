<script lang="ts" generics="T extends string | undefined">
	import type { Snippet } from 'svelte';
	import * as tooltip from '@zag-js/tooltip';
	import { useMachine, normalizeProps, type PropTypes, portal, mergeProps } from '@zag-js/svelte';
	import { nanoid } from 'nanoid';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { scale } from 'svelte/transition';

	type TooltipProps = {
		text?: string;
		placement?: tooltip.Placement;
		tag?: T;
		wrapper?: Snippet<[ReturnType<tooltip.Api<PropTypes>['getTriggerProps']>, Snippet]>;
		disable?: boolean;
		children?: Snippet;
	} & (T extends string ? SvelteHTMLElements[T] : SvelteHTMLElements['span']);

	const {
		text,
		placement = 'right',
		tag = 'span' as T,
		wrapper,
		disable,
		id,
		children,
		...props
	}: TooltipProps = $props();

	const service = useMachine(tooltip.machine, {
		id: id ?? nanoid(6),
		ids: { trigger: id },
		openDelay: 150,
		closeDelay: 150,
		positioning: {
			placement,
			strategy: 'fixed'
		}
	});
	const api: tooltip.Api<PropTypes> = $derived(
		!disable && text
			? tooltip.connect(service, normalizeProps)
			: {
					open: false,
					getArrowProps: () => ({}),
					getArrowTipProps: () => ({}),
					getContentProps: () => ({}),
					getPositionerProps: () => ({}),
					getTriggerProps: () => ({}),
					reposition: () => {},
					setOpen: () => {}
				}
	);
</script>

{#snippet content()}
	{#if api.open}
		<div use:portal {...api.getPositionerProps()}>
			<div
				transition:scale={{ duration: 150, start: 0.75, opacity: 0 }}
				class="z-50 rounded-lg border border-zinc-300 bg-zinc-200 px-2 py-1 text-sm/6 dark:border-zinc-700 dark:bg-zinc-800"
				{...api.getContentProps()}
			>
				{text}
			</div>
		</div>
	{/if}
{/snippet}

{#if wrapper}
	{@render wrapper(api.getTriggerProps(), content)}
{:else if tag}
	<svelte:element this={tag} {...mergeProps(api.getTriggerProps(), props)}>
		{@render content()}
		{@render children?.()}
	</svelte:element>
{/if}
