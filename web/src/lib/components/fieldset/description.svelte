<script lang="ts">
	import { clsx } from 'clsx';
	import { getContext } from 'svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import type { FieldContext } from './field.svelte';
	import { nanoid } from 'nanoid';

	type DescriptionProps = SvelteHTMLElements['p'];

	const { children, class: className, id, ...props }: DescriptionProps = $props();

	const { descriptionID } = getContext<FieldContext | undefined>('field') || {
		descriptionID: undefined
	};

	const randID = nanoid(6);

	$effect(() => {
		if (id) {
			descriptionID?.setValue(id);
		} else {
			descriptionID?.setValue(randID);
		}
	});
</script>

<p
	id={descriptionID?.value}
	data-slot="description"
	class={clsx(
		className,
		'data-disabled:opacity-50 text-base/6 text-zinc-500 sm:text-sm/6 dark:text-zinc-400'
	)}
	{...props}
>
	{#if children}
		{@render children()}
	{/if}
</p>
