<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SidebarHeaderProps = SvelteHTMLElements['div'];

	const { class: className, children, ...props }: SidebarHeaderProps = $props();
</script>

<div
	class={clsx(
		'flex flex-col border-b border-zinc-950/5 p-4 dark:border-white/5 [&>[data-slot=section]+[data-slot=section]]:mt-2.5',
		className
	)}
	{...props}
>
	{@render children?.()}
</div>
