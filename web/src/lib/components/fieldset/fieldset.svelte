<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type FieldsetProps = SvelteHTMLElements['fieldset'];

	const { children, class: className, ...props }: FieldsetProps = $props();
</script>

<fieldset
	class={clsx(className, '*:data-[slot=text]:mt-1 [&>*+[data-slot=control]]:mt-6')}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</fieldset>
