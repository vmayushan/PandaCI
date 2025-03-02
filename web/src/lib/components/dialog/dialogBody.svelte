<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DialogBodyProps = SvelteHTMLElements['div'];

	const { children, class: className, ...props }: DialogBodyProps = $props();
</script>

<div class={[className, 'mt-6']} {...props}>
	{#if children}
		{@render children()}
	{/if}
</div>
