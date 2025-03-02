<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import clsx from 'clsx';

	type InputGroupProps = SvelteHTMLElements['span'] & {};

	const { children, class: className }: InputGroupProps = $props();
</script>

<span
	data-slot="control"
	class={clsx(
		'relative isolate',
		'[&_input]:has-[[data-slot=icon]:first-child]:pl-10 [&_input]:has-[[data-slot=icon]:last-child]:pr-10 sm:[&_input]:has-[[data-slot=icon]:first-child]:pl-8 sm:[&_input]:has-[[data-slot=icon]:last-child]:pr-8',
		'[&>[data-slot=icon]]:pointer-events-none [&>[data-slot=icon]]:absolute [&>[data-slot=icon]]:top-3 [&>[data-slot=icon]]:z-10 [&>[data-slot=icon]]:size-5 sm:[&>[data-slot=icon]]:top-2.5 sm:[&>[data-slot=icon]]:size-4',
		'[&>[data-slot=icon]:first-child]:left-3 sm:[&>[data-slot=icon]:first-child]:left-2.5 [&>[data-slot=icon]:last-child]:right-3 sm:[&>[data-slot=icon]:last-child]:right-2.5',
		'[&>[data-slot=icon]]:text-zinc-500 dark:[&>[data-slot=icon]]:text-zinc-400',
		className
	)}
>
	{@render children?.()}
</span>
