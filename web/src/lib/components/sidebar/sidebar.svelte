<script lang="ts" module>
	export interface SidebarContext {
		collapsed: boolean;
	}
</script>

<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { setContext } from 'svelte';

	type SidebarProps = SvelteHTMLElements['nav'] & {
		collapsed?: boolean;
	};

	const { collapsed, class: className, children, ...props }: SidebarProps = $props();

	const context = $state({
		collapsed: collapsed ?? false
	});

	$effect(() => {
		context.collapsed = collapsed ?? false;
	});

	setContext<SidebarContext>('sidebar', context);
</script>

<nav class={[className, 'flex h-full min-h-0 flex-col']} {...props}>
	{#if children}
		{@render children()}
	{/if}
</nav>
