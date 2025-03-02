<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SidebarFooterProps = SvelteHTMLElements['div'];

	const { class: className, children, ...props }: SidebarFooterProps = $props();
</script>

<div
	class={clsx(
		'flex flex-col border-t border-zinc-950/5 px-4 py-2 dark:border-white/5 [&>[data-slot=section]+[data-slot=section]]:mt-2.5',
		className
	)}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</div>
