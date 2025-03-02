<script lang="ts">
	import { clsx } from 'clsx';
	import type { SvelteHTMLElements } from 'svelte/elements';

	const colors = {
		default: 'text-on-surface-variant',
		emphasis: 'text-on-surface',
		error: 'text-red-500 dark:text-red-400',
		success: 'text-green-500 dark:text-green-400'
	} as const;

	type TextProps = SvelteHTMLElements['p'] & {
		variant?: keyof typeof colors;
		size?: 'base' | 'sm';
	};

	const {
		children,
		class: className,
		size = 'base',
		variant = 'default',
		...props
	}: TextProps = $props();
</script>

<p
	data-slot="text"
	{...props}
	class={clsx(
		className,
		colors[variant],
		size === 'base' && 'text-base/6 sm:text-sm/6',
		size === 'sm' && 'text-sm/6 sm:text-xs/6'
	)}
>
	{#if children}
		{@render children()}
	{/if}
</p>
