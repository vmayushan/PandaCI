<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { getTableContext } from './table.svelte';
	import { getTableRowContext } from './tableRow.svelte';

	type TableCellProps = {} & SvelteHTMLElements['td'];

	const { class: className, children, ...props }: TableCellProps = $props();

	const { grid, bleed } = getTableContext();
	const { href, target, title } = getTableRowContext();

	let cellRef = $state<HTMLTableCellElement>();
</script>

<td
	{...props}
	bind:this={cellRef}
	class={[
		className,
		'first:pl-(--gutter,--spacing(2)) last:pr-(--gutter,--spacing(2)) border-b border-b-zinc-950/10 px-4 py-2 font-medium dark:border-b-white/10',
		grid && 'border-l border-l-zinc-950/5 first:border-l-0 dark:border-l-white/5',
		!bleed && 'sm:first:pl-1 sm:last:pr-1'
	]}
>
	{#if href}
		<a
			data-row-link
			{href}
			{target}
			aria-label={title}
			tabIndex={cellRef?.previousElementSibling === null ? 0 : -1}
			class="focus:outline-hidden absolute inset-0"
		></a>
	{/if}
	{@render children?.()}
</td>
