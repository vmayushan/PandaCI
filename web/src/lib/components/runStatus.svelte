<script lang="ts">
	import type { RunConclusion, RunStatus } from '$lib/api/organization';
	import clsx from 'clsx';
	import { ArrowsClockwise, CaretDoubleRight, Check, Hourglass, X } from 'phosphor-svelte';
	import type { SvelteHTMLElements } from 'svelte/elements';

	type RunStatusProps = SvelteHTMLElements['span'] & {
		status: RunStatus | undefined;
		conclusion?: RunConclusion;
		subtle?: boolean;
	};

	const { status, conclusion, subtle, class: className, ...props }: RunStatusProps = $props();

	const colors = {
		status: {
			running: 'bg-blue-600 text-white',
			completed: '',
			pending: 'bg-indigo-600 text-white',
			queued: '',
			skeleton: 'animate-pulse bg-black/10 dark:bg-white/5'
		} as Record<RunStatus | 'skeleton', string>,
		conclusion: {
			success: subtle ? 'bg-zinc-400 dark:bg-zinc-600 text-white' : 'bg-green-600 text-white',
			failure: 'bg-red-600 text-white',
			skipped: 'bg-gray-600 text-white'
		} as Record<RunConclusion, string>
	} as const;
</script>

<span
	class={clsx(
		'block shrink-0 rounded-full p-1 *:data-[slot=icon]:h-full *:data-[slot=icon]:w-full',
		status !== 'completed' && colors.status[status || 'skeleton'],
		status === 'completed' && conclusion && colors.conclusion[conclusion],
		className
	)}
	{...props}
>
	{#if status !== 'completed'}
		{#if status === 'running'}
			<ArrowsClockwise weight="bold" class="animate-[spin_2s_linear_infinite]" data-slot="icon" />
		{:else if status === 'pending'}
			<Hourglass weight="bold" data-slot="icon" />
		{/if}
	{:else if conclusion === 'success'}
		<Check weight="bold" data-slot="icon" />
	{:else if conclusion === 'failure'}
		<X weight="bold" data-slot="icon" />
	{:else if conclusion === 'skipped'}
		<CaretDoubleRight weight="bold" data-slot="icon" />
	{/if}
</span>
