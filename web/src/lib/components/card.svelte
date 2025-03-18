<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type CardProps = SvelteHTMLElements['card'] & {
		hideMobileRing?: boolean;
	};

	const { children, class: className, hideMobileRing, ...props }: CardProps = $props();
</script>

<div
	{...props}
	class={[
		className,
		'p-(--gutter) row-start-2 w-full min-w-0 rounded-lg bg-white shadow-sm ring-zinc-950/10 [--gutter:--spacing(8)] sm:rounded-2xl sm:ring-1 dark:bg-zinc-900 dark:ring-white/10 forced-colors:outline',
		!hideMobileRing && 'ring-1'
	]}
>
	{@render children?.()}
</div>
