<script lang="ts" module>
	import type { SvelteHTMLElements } from 'svelte/elements';

	type BaseButtonProps = {
		color?: keyof typeof colors;
	};

	export type ButtonProps<T extends { href?: string } = object> = T['href'] extends string
		? T & BaseButtonProps & SvelteHTMLElements['a']
		: T & BaseButtonProps & SvelteHTMLElements['button'];
</script>

<script lang="ts" generics="T extends {href?: string}">
	import TouchTarget from './touchTarget.svelte';
	import Badge, { colors } from './badge.svelte';

	const { children, class: className, color, ...rest }: ButtonProps<T> = $props();

	const classes = [
		className,
		'group relative inline-flex rounded-md focus:outline-hidden focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500'
	];

	const component = (rest.href ? 'a' : 'button') as 'a' | 'button';
</script>

<svelte:element this={component} class={classes} {...rest}>
	<Badge {color}>
		<TouchTarget />
		{@render children?.()}
	</Badge>
</svelte:element>
