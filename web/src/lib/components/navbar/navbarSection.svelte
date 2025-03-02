<script lang="ts" module>
	import { getContext } from 'svelte';

	export type NavbarSectionContext = {
		animationName?: string;
	};

	export function getNavbarSectionContext(): NavbarSectionContext {
		return getContext<NavbarSectionContext>('navbar-section');
	}
</script>

<script lang="ts">
	import { setContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type NavbarSectionProps = SvelteHTMLElements['div'];

	const id = $props.id();

	const { class: className, children, ...rest }: NavbarSectionProps = $props();

	setContext<NavbarSectionContext>('navbar-section', { animationName: `navbar-${id}` });
</script>

<div class={[className, 'flex items-center gap-3']} {...rest}>
	{@render children?.()}
</div>
