<script lang="ts">
	import type { StepRun } from '$lib/api/organization';
	import { CaretRight, Spinner } from 'phosphor-svelte';
	import { page } from '$app/state';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import RunLogs from '../../runLogs.svelte';
	import clsx from 'clsx';
	import { buttonStyles } from '$lib/components/button.svelte';
	import { Logs } from '../../workflow/logs.svelte';
	import RunStatus from '$lib/components/runStatus.svelte';
	import Text from '$lib/components/text/text.svelte';
	import LiveDate from '../../liveDate.svelte';

	interface StepLogProps {
		step: StepRun;
		workflowID: string;
	}

	const { step, workflowID }: StepLogProps = $props();

	let open = $state(step.status === 'completed' && step.conclusion === 'failure');

	let prev_id: string | undefined = '';

	const logMeta = createQuery(() => ({
		...queries.runs
			.projectRuns(page.params.orgName, page.params.projectName)
			._ctx.logs(workflowID, step.id),
		staleTime: Infinity,
		enabled: open && step.status === 'completed'
	}));

	// svelte-ignore state_referenced_locally
	const rawLogs = new Logs(page.params.orgName, page.params.projectName, 'ansi', open);

	$effect(() => {
		if (prev_id !== step.id) {
			open = step.status === 'completed' && step.conclusion === 'failure';
			prev_id = step.id;
		}

		rawLogs.update({
			logURL: logMeta.data?.url,
			runID: workflowID,
			status: step.status,
			stepID: step.id,
			enabled: open
		});
	});
</script>

<details bind:open class="group relative open:[&>summary>span>svg]:first:rotate-90">
	<summary class="bg-surface-low sticky top-0 flex">
		<span
			class={clsx(
				'group-open:bg-surface flex w-full items-center justify-start overflow-hidden',
				buttonStyles.base,
				buttonStyles.plain,
				buttonStyles.colors['dark/zinc']
			)}
		>
			{#if open && (rawLogs.isLoading || logMeta.isLoading)}
				<Spinner data-slot="icon" class="size-4 shrink-0 animate-spin" />
			{:else}
				<CaretRight data-slot="icon" class="size-4 shrink-0 transition-transform" />
			{/if}
			<RunStatus subtle conclusion={step.conclusion} status={step.status} class="mx-1.5 size-4" />
			<span
				class={clsx(
					'truncate whitespace-nowrap',
					step.conclusion === 'failure' && 'text-red-500 dark:text-red-400'
				)}
			>
				{step.name}
			</span>
			<Text class="ml-auto whitespace-nowrap" size="sm">
				<LiveDate startedAt={step.createdAt} finishedAt={step.finishedAt} />
			</Text>
		</span>
	</summary>

	<RunLogs class="pl-2">
		{#if rawLogs.htmlLogs}
			<div class="mt-2">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html rawLogs.htmlLogs}
			</div>
		{/if}
	</RunLogs>
</details>
