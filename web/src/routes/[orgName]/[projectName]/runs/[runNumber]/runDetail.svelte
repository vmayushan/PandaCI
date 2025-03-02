<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type RunDetailProps = SvelteHTMLElements['div'] & {
		label: string;
	};

	const { label, class: className, children, ...props }: RunDetailProps = $props();
</script>

<div class={clsx('space-y-2', className)} {...props}>
	<dt class="text-on-surface-variant text-sm">{label}</dt>
	<dd class="text-on-surface flex items-center space-x-2">{@render children?.()}</dd>
</div>
