<script lang="ts" module>
	import type { Api } from '@zag-js/combobox';
	export type ComboboxContext = {
		readonly api: Api;
	};
</script>

<script lang="ts" generics="T extends {label: string; onClick?: any; value: string; icon?: any;}">
	import * as combobox from '@zag-js/combobox';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { nanoid } from 'nanoid';
	import { setContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import Divider from '../divider.svelte';

	type ComboboxProps = SvelteHTMLElements['div'] & {
		items: T[];
		actions?: T[];
		itemOptions?: Omit<combobox.CollectionOptions<T>, 'items'>;
		selected?: T;
	};

	let {
		children,
		items,
		actions = [],
		selected = $bindable(),
		itemOptions,
		...props
	}: ComboboxProps = $props();

	let filtered = $state.raw([...items, ...actions]);

	const collection = $derived(
		combobox.collection({
			items: filtered,
			...itemOptions
		})
	);

	function createComboboxAPI() {
		const id = nanoid(6);

		const service = useMachine(combobox.machine, {
			id,
			collection,
			selectionBehavior: 'replace',
			loopFocus: true,
			defaultValue: selected ? [collection.getItemValue(selected)!] : undefined,
			allowCustomValue: true,
			multiple: false,
			openOnClick: true,
			onOpenChange() {
				collection.setItems([...items, ...actions]);
				filtered = [...items, ...actions];
				api.setInputValue((api.selectedItems[0] as any).label);
			},
			onInputValueChange({ inputValue }) {
				const newOptions =
					items.filter(
						(item) =>
							collection.stringifyItem(item)?.toLowerCase().includes(inputValue.toLowerCase()) ||
							Boolean(item.onClick)
					) || items;
				filtered = [...newOptions, ...actions];
				collection.setItems([...newOptions, ...actions]);
			},
			onValueChange({ items }: { items: T[] }) {
				const item: T = items[0];
				if (item?.onClick) {
					item.onClick();
					api.setValue([]);
				} else {
					selected = item;
				}
			}
		});

		const api = $derived(combobox.connect(service, normalizeProps));

		return {
			get api() {
				return api;
			}
		};
	}

	const ctx = createComboboxAPI();

	$effect(() => {
		ctx.api.collection.setItems(filtered);
	});

	setContext<ComboboxContext>('combobox', ctx);
</script>

<div data-slot="control" {...props}>
	{#if children}
		{@render children()}
	{/if}

	<div {...ctx.api.getPositionerProps()}>
		{#if filtered.length > 0}
			<ul
				{...ctx.api.getContentProps()}
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
				{#each filtered as item (item.label)}
					{#if actions[0]?.value === item.value}
						<Divider class="my-1" />
					{/if}
					<li
						{...ctx.api.getItemProps({ item })}
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
						<span
							class={[
								// Base
								'col-span-2 col-start-1 flex min-w-0 items-center',
								// Icons
								'*:data-[slot=icon]:size-5 *:data-[slot=icon]:shrink-0 sm:*:data-[slot=icon]:size-4',
								'group-data-highlighted/option:*:data-[slot=icon]:text-white *:data-[slot=icon]:text-zinc-500 dark:*:data-[slot=icon]:text-zinc-400',
								'forced-colors:group-data-highlighted/option:*:data-[slot=icon]:text-[Canvas] forced-colors:*:data-[slot=icon]:text-[CanvasText]',
								// Avatars
								'*:data-[slot=avatar]:-mx-0.5 *:data-[slot=avatar]:size-6 sm:*:data-[slot=avatar]:size-5'
							]}
						>
							{#if item.icon}
								<item.icon data-slot="icon" class="mr-2" />
							{/if}
							{item.label}
						</span>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
</div>
