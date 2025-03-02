<script lang="ts">
	import type { ButtonProps } from '../button.svelte';

	import Button from '../button.svelte';
	import { getContext } from 'svelte';
	import type { DialogContext } from './dialog.svelte';

	type DialogButtonProps = ButtonProps & { close?: boolean };

	const { children, ...props }: DialogButtonProps = $props();

	const ctx = getContext<DialogContext>('dialog');
</script>

<Button {...ctx.api.getCloseTriggerProps()} {...props}>
	{#if children}
		{@render children()}
	{/if}
</Button>
