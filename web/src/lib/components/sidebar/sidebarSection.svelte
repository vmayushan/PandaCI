<script lang="ts" module>
	export type SidebarSectionContext = {
		animationName?: string;
	};
</script>

<script lang="ts">
	import clsx from 'clsx';
	import { nanoid } from 'nanoid';
	import { setContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type SidebarSectionProps = SvelteHTMLElements['div'];

	const { class: className, children, ...props }: SidebarSectionProps = $props();

	const animationName = nanoid(6);

	setContext<SidebarSectionContext>('sidebar-section', { animationName });
</script>

<div data-slot="section" class={clsx('flex flex-col gap-0.5', className)} {...props}>
	{#if children}
		{@render children()}
	{/if}
</div>
