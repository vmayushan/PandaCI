<script lang="ts" generics="T">
	import * as select from '@zag-js/select';
	import { normalizeProps, portal, useMachine } from '@zag-js/svelte';
	import { nanoid } from 'nanoid';
	import type { Snippet } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { CaretUpDown } from 'phosphor-svelte';

	type ListboxButtonProps = Omit<SvelteHTMLElements['button'], 'color'> & {};

	type ListboxProps = SvelteHTMLElements['div'] & {
		button: Snippet<
			[
				ListboxButtonProps,
				T | undefined,
				{
					icon: Snippet;
				}
			]
		>;
		sameWidth?: boolean;
		selected?: T;
		multi?: boolean;
		items: T[];
		header?: Snippet;
		footer?: Snippet;
		item: Snippet<[T]>;
		placeholder?: Snippet;
		itemOptions?: Omit<select.CollectionOptions<T>, 'items'>;
	};

	let {
		selected = $bindable(),
		items,
		itemOptions,
		header,
		footer,
		button,
		item,
		sameWidth,
		multi
	}: ListboxProps = $props();

	const id = nanoid(6);

	let clickedOpen = $state(false);

	const collection = select.collection({
		items,
		...itemOptions
	});

	const service = useMachine(select.machine, {
		id,
		positioning: { sameWidth },
		collection,
		loopFocus: true,
		multiple: multi,
		value: selected ? [collection.getItemValue(selected)!] : undefined
	});

	const api = $derived(select.connect(service, normalizeProps));

	$effect.pre(() => {
		if (api.open && clickedOpen) service.send({ type: 'ITEM.POINTER_LEAVE' });
		clickedOpen = false;
	});

	// svelte-ignore state_referenced_locally
	let prevSelected = $state(api.selectedItems[0]);

	$effect(() => {
		api.collection.setItems(items);
	});

	$effect(() => {
		if (prevSelected === selected) {
			selected = api.selectedItems[0];
			prevSelected = api.selectedItems[0];
		} else if (prevSelected === api.selectedItems[0]) {
			api.selectValue(api.collection.getItemValue(selected)!);
			prevSelected = selected;
		}
	});

	const triggerProps = $derived(api.getTriggerProps());
</script>

{#snippet icon()}
	<span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
		<CaretUpDown
			data-slot="icon"
			class="group-data-disabled:text-zinc-600 size-5 text-zinc-500 sm:size-4 dark:text-zinc-400 forced-colors:text-[CanvasText]"
		/>
	</span>
{/snippet}

<div {...api.getRootProps()}>
	<div {...api.getControlProps()}>
		{@render button(
			{
				...triggerProps,
				onclick: (e) => {
					if (!api.open) clickedOpen = true;
					triggerProps.onclick?.(e);
				}
			},
			api.selectedItems[0],
			{ icon }
		)}
	</div>

	<div use:portal {...api.getPositionerProps()}>
		<div
			{...api.getContentProps()}
			class={[
				// Base styles
				'isolate z-50 select-none scroll-py-1 rounded-xl p-1',
				// Invisible border that is only visible in `forced-colors` mode for accessibility purposes
				'focus:outline-hidden outline outline-transparent',
				// Handle scrolling when menu won't fit in viewport
				'overflow-auto overscroll-contain',
				// Popover background
				'bg-white/75 backdrop-blur-xl dark:bg-zinc-800/75',
				// Shadows
				'shadow-lg ring-1 ring-zinc-950/10 dark:ring-inset dark:ring-white/10'
			]}
		>
			{#if header}
				{@render header()}
			{/if}
			<ul class="" onmouseleave={() => service.send({ type: 'ITEM.POINTER_LEAVE' })}>
				{#each items as itemData (itemData)}
					<li
						{...api.getItemProps({ item: itemData })}
						class={[
							// Basic layout
							'group/option grid cursor-default grid-cols-[--spacing(5)_1fr] items-baseline gap-x-2 rounded-lg py-2.5 pl-2 pr-3.5 sm:grid-cols-[--spacing(4)_1fr] sm:py-1.5 sm:pl-1.5 sm:pr-3',
							// Typography
							'text-base/6 text-zinc-950 sm:text-sm/6 dark:text-white forced-colors:text-[CanvasText]',
							// Focus
							'outline-hidden data-highlighted:bg-blue-500 data-highlighted:text-white',
							// Forced colors mode
							'forced-colors:data-highlighted:bg-[Highlight] forced-colors:data-highlighted:text-[HighlightText] forced-color-adjust-none',
							// Disabled
							'data-disabled:opacity-50'
						]}
					>
						<svg
							{...api.getItemIndicatorProps({ item: itemData }) as Record<string, unknown>}
							class="relative hidden size-5 self-center stroke-current group-aria-selected/option:inline sm:size-4"
							viewBox="0 0 16 16"
							fill="none"
							aria-hidden="true"
						>
							<path
								d="M4 8.5l3 3L12 4"
								stroke-width={1.5}
								stroke-linecap="round"
								stroke-linejoin="round"
							/>
						</svg>
						<span
							class={[
								// Base
								'col-start-2 flex min-w-0 items-center',
								// Icons
								'*:data-[slot=icon]:size-5 *:data-[slot=icon]:shrink-0 sm:*:data-[slot=icon]:size-4',
								'group-data-highlighted/option:*:data-[slot=icon]:text-white *:data-[slot=icon]:text-zinc-500 dark:*:data-[slot=icon]:text-zinc-400',
								'forced-colors:group-data-highlighted/option:*:data-[slot=icon]:text-[Canvas] forced-colors:*:data-[slot=icon]:text-[CanvasText]',
								// Avatars
								'*:data-[slot=avatar]:-mx-0.5 *:data-[slot=avatar]:size-6 sm:*:data-[slot=avatar]:size-5'
							]}
						>
							{@render item(itemData)}
						</span>
					</li>
				{/each}
			</ul>
			{#if footer}
				{@render footer()}
			{/if}
		</div>
	</div>
</div>
