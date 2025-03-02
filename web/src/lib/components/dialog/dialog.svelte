<script lang="ts" module>
	import type { Api } from '@zag-js/dialog';
	export type DialogContext = {
		readonly api: Api;
	};
</script>

<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import * as dialog from '@zag-js/dialog';
	import { portal, normalizeProps, useMachine } from '@zag-js/svelte';
	import Card from '../card.svelte';
	import { nanoid } from 'nanoid';
	import { setContext } from 'svelte';

	const sizes = {
		xs: 'sm:max-w-xs',
		sm: 'sm:max-w-sm',
		md: 'sm:max-w-md',
		lg: 'sm:max-w-lg',
		xl: 'sm:max-w-xl',
		'2xl': 'sm:max-w-2xl',
		'3xl': 'sm:max-w-3xl',
		'4xl': 'sm:max-w-4xl',
		'5xl': 'sm:max-w-5xl'
	};

	type DialogProps = SvelteHTMLElements['div'] & {
		size?: keyof typeof sizes;
		open?: boolean;
	};

	let {
		children,
		class: className,
		size = 'lg',
		open = $bindable(),
		...props
	}: DialogProps = $props();

	const id = nanoid(6);

	const service = useMachine(dialog.machine, { id });

	const api = $derived(dialog.connect(service, normalizeProps));

	let prevOpen = $state(false);

	$effect(() => {
		if (open === undefined) {
			prevOpen = api.open;
			return;
		}

		if (prevOpen === open) {
			open = api.open;
			prevOpen = api.open;
		} else {
			api.setOpen(open);
			prevOpen = open;
		}
	});

	setContext<DialogContext>('dialog', {
		get api() {
			return api;
		}
	});
</script>

{#if api.open}
	<div
		use:portal
		{...api.getBackdropProps()}
		class="fixed inset-0 flex w-screen justify-center overflow-y-auto bg-zinc-950/25 px-2 py-2 focus:outline-0 sm:px-6 sm:py-8 lg:px-8 lg:py-16 dark:bg-zinc-950/50"
	></div>
	<div use:portal {...api.getPositionerProps()} class="fixed inset-0 z-10 w-screen overflow-y-auto">
		<div class="flex min-h-full items-center justify-center p-4">
			<Card {...api.getContentProps()} class={['h-full', sizes[size], className]} {...props}>
				{#if children}
					{@render children()}
				{/if}
			</Card>
		</div>
	</div>
{/if}
