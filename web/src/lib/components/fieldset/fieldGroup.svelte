<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type FieldGroupProps = SvelteHTMLElements['div'];

	const { children, class: className, ...props }: FieldGroupProps = $props();
</script>

<div data-slot="control" class={clsx(className, 'space-y-8')} {...props}>
	{#if children}
		{@render children()}
	{/if}
</div>
