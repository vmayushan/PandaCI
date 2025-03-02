<script lang="ts">
	import * as avatar from '@zag-js/avatar';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { nanoid } from 'nanoid';
	import type { SvelteHTMLElements } from 'svelte/elements';

	function getInitials(name: string) {
		const words = name.split(' ');
		if (words.length <= 1) {
			return words[0]?.[0]?.toUpperCase(); // Only one name, return its initial
		}
		return (words[0][0] + words[words.length - 1][0]).toUpperCase();
	}

	type AvatarProps = SvelteHTMLElements['span'] & {
		src?: string;
		name?: string;
		square?: boolean;
		alt: string;
		class?: string;
	};

	const { src, children, name, class: className, alt, square, ...props }: AvatarProps = $props();

	const id = nanoid(6);

	const service = useMachine(avatar.machine, { id });
	const api = $derived(avatar.connect(service, normalizeProps));
</script>

<span
	data-slot="avatar"
	{...props}
	{...api.getRootProps()}
	class={[
		className,
		// Basic layout
		'inline-grid shrink-0 align-middle [--avatar-radius:20%] [--ring-opacity:20%] *:col-start-1 *:row-start-1',
		'outline-black/(--ring-opacity) dark:outline-white/(--ring-opacity) outline -outline-offset-1',
		// Add the correct border radius
		square ? 'rounded-(--avatar-radius) *:rounded-(--avatar-radius)' : 'rounded-full *:rounded-full'
	]}
>
	{#if name}
		<svg
			class="size-full select-none fill-current p-[5%] text-[48px] font-medium uppercase"
			viewBox="0 0 100 100"
			{...api.getFallbackProps() as Record<string, unknown>}
		>
			<title>{alt}</title>
			<text
				x="50%"
				y="50%"
				alignment-baseline="middle"
				dominant-baseline="middle"
				text-anchor="middle"
				dy=".125em"
			>
				{getInitials(name)}
			</text>
		</svg>
	{/if}

	{@render children?.()}

	{#if src}
		<img class="size-full object-cover" {src} {alt} {...api.getImageProps()} />
	{/if}
</span>
