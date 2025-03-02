<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type CardProps = SvelteHTMLElements['card'];

	const { children, class: className, ...props }: CardProps = $props();
</script>

<div
	{...props}
	class={[
		className,
		'p-(--gutter) row-start-2 w-full min-w-0 rounded-lg bg-white shadow-sm ring-1 ring-zinc-950/10 [--gutter:--spacing(8)] sm:rounded-2xl dark:bg-zinc-900 dark:ring-white/10 forced-colors:outline'
	]}
>
	{@render children?.()}
</div>
