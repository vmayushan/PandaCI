<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type TableHeadProps = {} & SvelteHTMLElements[' thead'];

	const { children, class: className, ...props }: TableHeadProps = $props();
</script>

<thead {...props} class={[className, 'text-zinc-500 dark:text-zinc-400']}>
	{@render children?.()}
</thead>
