<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SidebarBodyProps = SvelteHTMLElements['div'];

	const { class: className, children, ...props }: SidebarBodyProps = $props();
</script>

<div
	class={clsx(
		'flex flex-1 flex-col overflow-y-auto overflow-x-hidden p-4 [&>[data-slot=section]+[data-slot=section]]:mt-8',
		className
	)}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</div>
