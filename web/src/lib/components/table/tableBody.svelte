<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type TableBodyProps = {} & SvelteHTMLElements[' thead'];

	const { children, ...props }: TableBodyProps = $props();
</script>

<tbody {...props}>
	{@render children?.()}
</tbody>
