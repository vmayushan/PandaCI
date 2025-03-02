<script lang="ts">
	import clsx from 'clsx';
	import { getContext } from 'svelte';
	import type { HTMLLabelAttributes } from 'svelte/elements';
	import type { FieldContext } from './field.svelte';
	import { nanoid } from 'nanoid';

	type Props = HTMLLabelAttributes;
	const { class: className, id, children, ...props }: Props = $props();

	const { inputID, labelID } = getContext<FieldContext | undefined>('field') || {
		inputID: undefined,
		labelID: undefined
	};

	const randID = nanoid(6);

	$effect(() => {
		if (id) {
			labelID?.setValue(id);
		} else {
			labelID?.setValue(randID);
		}
	});
</script>

<label
	data-slot="label"
	for={inputID?.value}
	id={id || randID}
	{...props}
	class={clsx(
		className,
		'data-disabled:opacity-50 select-none text-base/6 text-zinc-950 sm:text-sm/6 dark:text-white'
	)}
>
	{#if children}
		{@render children()}
	{/if}
</label>
