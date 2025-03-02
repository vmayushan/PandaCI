<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { getTableContext } from './table.svelte';

	type TableHeaderProps = {} & SvelteHTMLElements['th'];

	const { class: className, children, ...props }: TableHeaderProps = $props();

	const { grid, bleed } = getTableContext();
</script>

<th
	{...props}
	class={[
		className,
		'first:pl-(--gutter,--spacing(2)) last:pr-(--gutter,--spacing(2)) border-b border-b-zinc-950/10 px-4 py-2 font-medium dark:border-b-white/10',
		grid && 'border-l border-l-zinc-950/5 first:border-l-0 dark:border-l-white/5',
		!bleed && 'sm:first:pl-1 sm:last:pr-1'
	]}
>
	{@render children?.()}
</th>
