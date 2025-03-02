<script lang="ts">
	import type { ButtonProps } from '../button.svelte';
	import Button from '../button.svelte';
	import { getContext, type Snippet } from 'svelte';
	import type { DropdownContext } from './dropdown.svelte';

	type DropdownButtonProps = ButtonProps & {
		indicator?: Snippet<[Record<string, unknown>]>;
		renderer?: Snippet<[Record<string, unknown>, Record<string, unknown>]>;
	};

	const { children, indicator, renderer, ...props }: DropdownButtonProps = $props();

	const ctx = getContext<DropdownContext>('dropdown');
</script>

{#snippet defaultRenderer(attr: any, indicatorProps: any)}
	<Button {...attr}>
		{#if children}
			{@render children()}
		{/if}
		{#if indicator}
			{@render indicator(indicatorProps)}
		{/if}
	</Button>
{/snippet}

{@render (renderer ?? defaultRenderer)(
	{ ...ctx.api.getTriggerProps(), ...props },
	{
		...ctx.api.getIndicatorProps(),
		'data-slot': 'icon'
	}
)}
