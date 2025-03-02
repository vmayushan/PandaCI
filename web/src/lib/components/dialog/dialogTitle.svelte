<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DialogTitleProps = SvelteHTMLElements['h2'];

	const { children, class: className, ...props }: DialogTitleProps = $props();
</script>

<h2
	class={[
		className,
		'text-balance text-lg/6 font-semibold text-zinc-950 sm:text-base/6 dark:text-white'
	]}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</h2>
