<script lang="ts">
	import clsx from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DividerProps = SvelteHTMLElements['hr'] & { soft?: boolean };

	const { class: className, soft, ...props }: DividerProps = $props();
</script>

<hr
	role="presentation"
	{...props}
	class={clsx(
		className,
		'w-full border-t',
		soft && 'border-zinc-950/5 dark:border-white/5',
		!soft && 'border-zinc-950/10 dark:border-white/10'
	)}
/>
