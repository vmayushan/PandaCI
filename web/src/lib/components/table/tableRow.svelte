<script lang="ts" module>
	import { getContext } from 'svelte';
	export type TableRowContext = {
		readonly href?: string;
		readonly target?: string;
		readonly title?: string;
	};

	export function getTableRowContext(): TableRowContext {
		return getContext<TableRowContext>('table-row');
	}
</script>

<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { setContext } from 'svelte';
	import { getTableContext } from './table.svelte';

	type TableRowProps = {
		href?: string;
		target?: string;
		title?: string;
	} & SvelteHTMLElements['tr'];

	const { class: className, children, href, target, title, ...props }: TableRowProps = $props();

	const { striped } = getTableContext();

	setContext<TableRowContext>('table-row', { href, target, title });
</script>

<tr
	{...props}
	class={[
		className,
		href &&
			'relative has-[[data-row-link]:focus]:outline-2 has-[[data-row-link]:focus]:-outline-offset-2 has-[[data-row-link]:focus]:outline-blue-500 dark:focus-within:bg-white/[2.5%]',
		striped && 'even:bg-zinc-950/[2.5%] dark:even:bg-white/[2.5%]',
		href && striped && 'hover:bg-zinc-950/5 dark:hover:bg-white/5',
		href && !striped && 'hover:bg-zinc-950/[2.5%] dark:hover:bg-white/[2.5%]'
	]}
>
	{@render children?.()}
</tr>
