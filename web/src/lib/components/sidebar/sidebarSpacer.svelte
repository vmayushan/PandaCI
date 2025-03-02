<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SidebarSpacerProps = Omit<SvelteHTMLElements['div'], 'children'>;

	const { class: className, ...props }: SidebarSpacerProps = $props();
</script>

<div aria-hidden="true" class={clsx('mt-8 flex-1', className)} {...props}></div>
