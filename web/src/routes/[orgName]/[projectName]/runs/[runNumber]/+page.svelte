<script lang="ts">
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import { Skeleton } from '$lib/components';
	import LiveDate from './liveDate.svelte';
	import SubHeading from '$lib/components/subHeading.svelte';
	import Alert from './alert.svelte';
	import RunDetail from './runDetail.svelte';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import Tooltip from '$lib/components/tooltip/tooltip.svelte';
	import Avatar from '$lib/components/avatar.svelte';

	dayjs.extend(relativeDate);

	const run = createQuery(() =>
		queries.runs
			.projectRuns(page.params.orgName, page.params.projectName)
			._ctx.get(page.params.runNumber)
	);
</script>

<div class="flex h-min w-full flex-col space-y-4">
	<div class="border-outline-variant bg-surface h-min w-full overflow-auto rounded-lg border p-4">
		<SubHeading>Workflow Summary</SubHeading>

		<dl class="my-4 grid grid-cols-1 gap-4 lg:grid-cols-2">
			<RunDetail label="Duration">
				{#if run.data}
					<LiveDate startedAt={run.data.createdAt} finishedAt={run.data.finishedAt} />
				{:else}
					<Skeleton class="mt-4 h-4 w-20" />
				{/if}
			</RunDetail>
			<RunDetail label="Runner">
				{#if run.data}
					ubuntu-2x
				{:else}
					<Skeleton class="mt-4 h-4 w-44" />
				{/if}
			</RunDetail>
			<RunDetail label="Branch">
				{#if run.data}
					{run.data.gitBranch}
				{:else}
					<Skeleton class="mt-4 h-4 w-44" />
				{/if}
			</RunDetail>
			<RunDetail label="Commit">
				{#if run.data}
					{run.data.gitSha}
				{:else}
					<Skeleton class="mt-4 h-4 w-96" />
				{/if}
			</RunDetail>
			<RunDetail label="Trigger">
				{#if run.data}
					{run.data.trigger}
				{:else}
					<Skeleton class="mt-4 h-4 w-96" />
				{/if}
			</RunDetail>
			<RunDetail label="Pr">
				{#if run.data}
					{run.data.prNumber ? `#${run.data.prNumber}` : 'N/A'}
				{:else}
					<Skeleton class="mt-4 h-4 w-24" />
				{/if}
			</RunDetail>
			<RunDetail label="Committer email">
				{#if run.data}
					{#if run.data.committer.avatar}
						<Avatar class="size-5" src={run.data.committer.avatar} alt={run.data.committer.email} />
					{/if}
					<span>{run.data.committer.email}</span>
				{:else}
					<Skeleton class="mt-4 h-4 w-48" />
				{/if}
			</RunDetail>
			<RunDetail label="Committer name">
				{#if run.data}
					{run.data.committer.name ?? 'N/A'}
				{:else}
					<Skeleton class="mt-4 h-4 w-36" />
				{/if}
			</RunDetail>
			<RunDetail label="Created">
				{#if run.data}
					<Tooltip text={run.data.createdAt}>
						{dayjs(run.data.createdAt).fromNow()}
					</Tooltip>
				{:else}
					<Skeleton class="mt-4 h-4 w-32" />
				{/if}
			</RunDetail>
			<RunDetail label="Finished">
				{#if run.data}
					{#if run.data.finishedAt}
						<Tooltip text={run.data.finishedAt}>
							{dayjs(run.data.finishedAt).fromNow()}
						</Tooltip>
					{:else}
						Still running
					{/if}
				{:else}
					<Skeleton class="mt-4 h-4 w-32" />
				{/if}
			</RunDetail>
		</dl>
	</div>

	{#if run.data?.alerts?.length}
		<div class="border-outline-variant bg-surface h-min w-full overflow-auto rounded-lg border p-4">
			<SubHeading>Alerts</SubHeading>
			<ul class="mt-4 flex flex-col space-y-2">
				{#each run.data?.alerts ?? [] as alert (alert)}
					<li>
						<Alert title={alert.title} message={alert.message} type={alert.type} />
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</div>
