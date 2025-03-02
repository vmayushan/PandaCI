<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DescriptionListProps = SvelteHTMLElements['dl'] & {};

	const { children, class: className, ...props }: DescriptionListProps = $props();
</script>

<dl
	class={clsx(
		'grid grid-cols-1 text-base/6 sm:grid-cols-[min(50%,theme(spacing.80))_auto] sm:text-sm/6',
		className
	)}
	{...props}
>
	{@render children?.()}
</dl>
