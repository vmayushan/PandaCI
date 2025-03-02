<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DropdownHeadingProps = SvelteHTMLElements['header'];

	const { children, class: className, ...props }: DropdownHeadingProps = $props();
</script>

<div {...props} data-slot="label" class={[className, 'col-start-2 row-start-1']}>
	{@render children?.()}
</div>
