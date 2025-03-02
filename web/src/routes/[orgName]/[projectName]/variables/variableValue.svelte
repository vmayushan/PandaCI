<script lang="ts">
	import { queries } from '$lib/queries';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { Button } from '$lib/components';
	import { Eye, EyeSlash, Lock } from 'phosphor-svelte';
	import { buttonStyles } from '$lib/components/button.svelte';
	import Tooltip from '$lib/components/tooltip/tooltip.svelte';

	interface Props {
		id: string;
		sensitive: boolean;
	}

	const { id, sensitive }: Props = $props();

	let show = $state(false);

	const variable = createQuery(() => ({
		...queries.variables
			.projectVariables(page.params.orgName, page.params.projectName)
			._ctx.get(id),
		enabled: show
	}));
</script>

<div class="flex items-center space-x-2">
	{#if sensitive}
		<Tooltip text="Sensitive">
			<span
				class={[buttonStyles.base, buttonStyles.plain, buttonStyles.inner, 'hover:!bg-transparent']}
			>
				<Lock data-slot="icon" />
			</span>
		</Tooltip>
	{:else}
		<Button
			loading={variable.isLoading}
			plain
			disabled={variable.isLoading}
			onclick={() => (show = !show)}
		>
			{#if !show}
				<Eye data-slot="icon" />
			{:else}
				<EyeSlash data-slot="icon" />
			{/if}
		</Button>
	{/if}
	<span class={show && !variable.isLoading ? 'text-on-surface' : 'text-on-surface-variant'}>
		{(show ? variable.data?.value : '') || '•••••••••••••••'}
	</span>
</div>
