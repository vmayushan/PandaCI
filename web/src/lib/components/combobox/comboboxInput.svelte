<script lang="ts">
	import { getContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import type { ComboboxContext } from './combobox.svelte';
	import { CaretDown } from 'phosphor-svelte';
	import Button from '../button.svelte';
	import type { FieldContext } from '../fieldset/field.svelte';

	type ComboboxInputProps = Omit<HTMLInputAttributes, 'id'>;

	let { class: className, ...props }: ComboboxInputProps = $props();

	const { api } = getContext<ComboboxContext>('combobox');

	const { inputID, labelID } = getContext<FieldContext | undefined>('field') ?? {};

	$effect(() => {
		inputID?.setValue(api.getInputProps().id);
	});
</script>

<span
	{...api.getControlProps()}
	class={[
		className,
		// Basic layout
		'relative flex w-full rounded-lg',
		// Background color + shadow applied to inset pseudo element, so shadow blends with border in light mode
		'before:absolute before:inset-px before:rounded-[calc(var(--radius-lg)-1px)] before:bg-white before:shadow-sm',
		// Background color is moved to control and shadow is removed in dark mode so hide `before` pseudo
		'dark:before:hidden',
		// Focus ring
		'after:pointer-events-none after:absolute after:inset-0 after:rounded-lg after:ring-inset after:ring-transparent sm:focus-within:after:ring-2 sm:focus-within:after:ring-blue-500',
		// Disabled state
		'has-disabled:opacity-50 has-disabled:before:bg-zinc-950/5 has-disabled:before:shadow-none',
		// user-invalid state
		'has-[:user-invalid]:before:shadow-red-500/10',
		// Background color
		'bg-transparent dark:bg-white/5',
		// Border
		'border border-zinc-950/10 hover:border-zinc-950/20 dark:border-white/10 dark:hover:border-white/20'
	]}
>
	<input
		{...api.getInputProps()}
		aria-labelledby={labelID?.value}
		class={[
			// Basic layout
			'relative block w-full appearance-none px-[calc(--spacing(3.5)-1px)] py-[calc(--spacing(2.5)-1px)] sm:px-[calc(--spacing(3)-1px)] sm:py-[calc(--spacing(1.5)-1px)]',
			// Typography
			'text-base/6 text-zinc-950 placeholder:text-zinc-500 sm:text-sm/6 dark:text-white',
			// Background color
			'bg-transparent',
			// Hide default focus styles
			'focus:outline-hidden',
			// [&:user-invalid] state
			'[&:user-invalid]:border-red-500 [&:user-invalid]:hover:border-red-500 dark:[&:user-invalid]:border-red-500 dark:[&:user-invalid]:hover:border-red-500',
			// Disabled state
			'disabled:border-zinc-950/20 dark:disabled:border-white/15 dark:disabled:bg-white/[2.5%] dark:hover:disabled:border-white/15',
			// System icons
			'dark:[color-scheme:dark]'
		]}
		{...props}
	/>
	<Button plain {...api.getTriggerProps()}><CaretDown /></Button>
</span>
