<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DialogActionsProps = SvelteHTMLElements['div'];

	const { children, class: className, ...props }: DialogActionsProps = $props();
</script>

<div
	class={[
		className,
		'mt-8 flex flex-col-reverse items-center justify-end gap-3 *:w-full sm:flex-row sm:*:w-auto'
	]}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</div>
