<script lang="ts">
	import clsx from 'clsx';
	import { type Snippet } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type DropdownSectionProps = SvelteHTMLElements['div'] & {
		children: Snippet;
	};

	const { children, class: className, ...props }: DropdownSectionProps = $props();
</script>

<div
	class={clsx(
		className,
		// Define grid at the section level instead of the item level if subgrid is supported
		'col-span-full w-full supports-[grid-template-columns:subgrid]:grid supports-[grid-template-columns:subgrid]:grid-cols-[auto_1fr_1.5rem_0.5rem_auto]'
	)}
	role="group"
	{...props}
>
	{@render children()}
</div>
