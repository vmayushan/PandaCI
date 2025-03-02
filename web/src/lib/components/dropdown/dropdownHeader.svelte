<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DropdownHeaderProps = SvelteHTMLElements['div'];

	const { children, class: className, ...props }: DropdownHeaderProps = $props();
</script>

<div {...props} class={clsx(className, 'col-span-5 px-3.5 pb-1 pt-2.5 sm:px-3')}>
	{@render children?.()}
</div>
