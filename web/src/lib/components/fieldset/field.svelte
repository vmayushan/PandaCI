<script lang="ts" module>
	export interface FieldContext {
		inputID: {
			readonly value: string;
			setValue: (value: string) => string;
		};
		labelID: {
			readonly value: string;
			setValue: (value: string) => string;
		};
		descriptionID: {
			readonly value: string;
			setValue: (value: string) => string;
		};
	}
</script>

<script lang="ts">
	import clsx from 'clsx';
	import type { HTMLAttributes } from 'svelte/elements';
	import { setContext } from 'svelte';

	type Props = HTMLAttributes<HTMLDivElement>;

	const { class: className, children, ...props }: Props = $props();

	let inputID = $state<string>('');
	let labelID = $state<string>('');
	let descriptionID = $state<string>('');

	const ctx = {
		inputID: {
			get value() {
				return inputID;
			},
			setValue: (value: string) => (inputID = value)
		},
		labelID: {
			get value() {
				return labelID;
			},
			setValue: (value: string) => (labelID = value)
		},
		descriptionID: {
			get value() {
				return descriptionID;
			},
			setValue: (value: string) => (descriptionID = value)
		}
	};

	setContext<FieldContext>('field', ctx);
</script>

<div
	{...props}
	class={clsx(
		className,
		'[&>[data-slot=label]+[data-slot=control]]:mt-3',
		'[&>[data-slot=label]+[data-slot=description]]:mt-1',
		'[&>[data-slot=description]+[data-slot=control]]:mt-3',
		'[&>[data-slot=control]+[data-slot=description]]:mt-3',
		'[&>[data-slot=control]+[data-slot=error]]:mt-3',
		'*:data-[slot=label]:font-medium'
	)}
>
	{#if children}
		{@render children()}
	{/if}
</div>
