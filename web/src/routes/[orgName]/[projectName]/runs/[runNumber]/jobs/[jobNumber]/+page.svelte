<script lang="ts">
	import { page } from '$app/state';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { queries } from '$lib/queries';
	import type { TaskRun, WorkflowRun } from '$lib/api/organization';
	import { Text, Heading, Badge, RunStatus } from '$lib/components';
	import StepLog from './stepLog.svelte';
	import LiveDate from '../../liveDate.svelte';

	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import Skeleton from '$lib/components/skeleton.svelte';
	import BadgeButton from '$lib/components/badgeButton.svelte';

	dayjs.extend(relativeDate);

	const queryClient = useQueryClient();

	const workflow = createQuery(() => ({
		...queries.runs
			.projectRuns(page.params.orgName, page.params.projectName)
			._ctx.get(page.params.runNumber),
		refetchInterval() {
			return queryClient.getQueryData<WorkflowRun | undefined>(
				queries.runs
					.projectRuns(page.params.orgName, page.params.projectName)
					._ctx.get(page.params.runNumber).queryKey
			)?.status !== 'completed'
				? 5000
				: false;
		}
	}));

	const job = $derived(
		workflow.data?.jobs?.find((job) => job.number.toString() === page.params.jobNumber)
	);

	const getDockerHubLink = (image?: string) => {
		const parts = image?.split(':') ?? [];
		const imageName = parts.slice(0, -1).join(':');

		// Not a docker hub registry
		if (!imageName || imageName?.includes('.')) return;

		// Check if it's an official image (no namespace specified)
		const isOfficial = !imageName.includes('/');

		return `https://hub.docker.com/${isOfficial ? '_' : 'r'}/${imageName}`;
	};
</script>

{#snippet Task(task: TaskRun)}
	{@const dockerHubLink = getDockerHubLink(task.dockerImage)}
	<div class="border-outline-variant bg-surface-low mt-4 h-min w-full rounded-lg border">
		<header
			class={[
				'border-outline-variant bg-surface flex items-center justify-between px-4 py-2',
				task.conclusion === 'skipped' ? 'rounded-lg' : 'rounded-t-lg border-b'
			]}
		>
			<Heading level={3} size="xs" class="flex items-center space-x-2">
				<RunStatus class="size-4" conclusion={task.conclusion} status={task.status} />
				<span>{task.name}</span>
			</Heading>

			<div class="flex items-center space-x-2">
				<Text size="sm">
					<LiveDate startedAt={task.createdAt} finishedAt={task.finishedAt} />
				</Text>
				<BadgeButton target={dockerHubLink ? '_blank' : undefined} href={dockerHubLink}>
					{task.dockerImage}
				</BadgeButton>
			</div>
		</header>
		{#if task.conclusion !== 'skipped'}
			<ul class="flex flex-col p-4">
				{#if task.steps.length === 0 && task.status !== 'completed'}
					<Skeleton class="mx-2 h-6 grow py-1.5" />
				{:else if task.steps.length === 0}
					<Text class="text-center">No steps were run</Text>
				{/if}

				{#each task.steps ?? [] as step (step.id)}
					<li>
						<StepLog {step} workflowID={workflow.data!.id} />
					</li>
				{/each}
			</ul>
		{/if}
	</div>
{/snippet}

<div class="flex min-w-0 grow flex-col">
	<header class="border-outline-variant flex justify-between rounded-lg border px-4 py-2">
		<div>
			<Heading class="flex items-center space-x-2" size="sm" level={2}>
				<RunStatus class="size-5" conclusion={job?.conclusion} status={job?.status} />
				<span>{job?.name}</span>
				{#if workflow.isLoading}
					<Skeleton class="my-1 h-6 w-44" />
				{/if}
			</Heading>
			<div class="flex items-center space-x-2">
				{#if job}
					<Text size="sm">
						{#if job.status === 'completed'}
							Job ran {dayjs(job.createdAt).fromNow()} in <LiveDate
								startedAt={job.createdAt}
								finishedAt={job.finishedAt}
							/>
						{:else}
							<LiveDate startedAt={job.createdAt} />
						{/if}
					</Text>
				{/if}
				{#if workflow.isLoading}
					<Skeleton class="my-0.5 h-5 w-44" />
				{/if}
			</div>
		</div>
		<div class="flex items-start justify-center space-x-2">
			{#if workflow.isLoading}
				<Skeleton class="h-6 w-20" />
			{/if}
			{#if job?.runner}
				<Badge>
					{job?.runner}
				</Badge>
			{/if}
		</div>
	</header>

	{#if workflow.isLoading}
		<Skeleton class="mt-4 h-32" />
		<Skeleton class="mt-4 h-32" />
	{/if}

	{#each job?.tasks ?? [] as task (task.id)}
		{@render Task(task)}
	{/each}
</div>
