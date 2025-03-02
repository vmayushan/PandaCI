<script lang="ts" module>
	import { getContext } from 'svelte';
	export type TableContext = {
		readonly bleed?: boolean;
		readonly dense?: boolean;
		readonly grid?: boolean;
		readonly striped?: boolean;
	};

	export function getTableContext(): TableContext {
		return getContext<TableContext>('table');
	}
</script>

<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { setContext } from 'svelte';

	type TableProps = {
		bleed?: boolean;
		dense?: boolean;
		grid?: boolean;
		striped?: boolean;
	} & SvelteHTMLElements['div'];

	const {
		class: className,
		bleed,
		children,
		dense,
		grid,
		striped,
		...props
	}: TableProps = $props();

	setContext<TableContext>('table', { bleed, dense, grid, striped });
</script>

<div class="flow-root">
	<div
		{...props}
		class={[className, 'scrollbar-thin -mx-(--gutter) overflow-x-auto whitespace-nowrap']}
	>
		<div class={['inline-block min-w-full align-middle', !bleed && 'sm:px-(--gutter)']}>
			<table class="min-w-full text-left text-sm/6 text-zinc-950 dark:text-white">
				{@render children?.()}
			</table>
		</div>
	</div>
</div>
