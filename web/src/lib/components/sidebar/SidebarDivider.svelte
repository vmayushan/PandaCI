<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SidebarDividerProps = SvelteHTMLElements['hr'];

	const { class: className, ...props }: SidebarDividerProps = $props();
</script>

<hr
	class={clsx('-mx-4 my-4 border-t border-zinc-950/5 dark:border-white/5', className)}
	{...props}
/>
