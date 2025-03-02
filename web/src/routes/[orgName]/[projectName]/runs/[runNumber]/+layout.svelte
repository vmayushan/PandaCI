<script lang="ts">
	import { page } from '$app/state';
	import { Title, Text } from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { CaretLeft, House, Notepad } from 'phosphor-svelte';
	import RunStatus from '$lib/components/runStatus.svelte';
	import type { WorkflowRun } from '$lib/api/organization';
	import SidebarItem from '$lib/components/sidebar/sidebarItem.svelte';
	import SidebarBody from '$lib/components/sidebar/sidebarBody.svelte';
	import SidebarDivider from '$lib/components/sidebar/SidebarDivider.svelte';
	import SidebarSection from '$lib/components/sidebar/sidebarSection.svelte';
	import SidebarHeading from '$lib/components/sidebar/sidebarHeading.svelte';
	import SidebarLabel from '$lib/components/sidebar/sidebarLabel.svelte';
	import type { Snippet } from 'svelte';
	import Sidebar from '$lib/components/sidebar/sidebar.svelte';
	import Skeleton from '$lib/components/skeleton.svelte';
	import LiveDate from './liveDate.svelte';

	const queryClient = useQueryClient();

	const run = createQuery(() => ({
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

	const { children }: { children: Snippet } = $props();

	const sortedJobs = $derived(
		run.data?.jobs
			?.slice()
			.sort((a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()) ?? []
	);

	const baseHref = `/${page.params.orgName}/${page.params.projectName}/runs/${page.params.runNumber}`;
</script>

<Title
	class="h-28"
	titleLoading={run.isLoading}
	state={{
		status: run.data?.status,
		conclusion: run.data?.conclusion
	}}
	title={run.data?.gitTitle ?? `#${run.data?.number}`}
>
	{#snippet description()}
		{#if run.isLoading}
			<Skeleton class="ml-10 mt-4 h-3.5 w-56" />
		{:else}
			<span class="ml-10 text-xs">{run.data?.name}</span>
		{/if}
	{/snippet}

	{#snippet action()}
		<div class="flex w-full justify-between">
			<a
				href={`/${page.params.orgName}/${page.params.projectName}`}
				class="focus-visible:text-on-surface focus-visible:ring-offset-surface text-on-surface-variant hover:text-on-surface focus:outline-hidden flex items-center space-x-1 rounded text-sm focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
			>
				<CaretLeft class="size-4" data-slot="icon" />
				<span>Back to runs</span>
			</a>
		</div>
	{/snippet}
</Title>

<div class="mt-4 flex h-fit space-x-4">
	<Sidebar class="sticky top-0 flex w-full max-w-60 shrink-0 grow flex-col 2xl:max-w-xs">
		<SidebarBody>
			<SidebarItem
				href={baseHref}
				current={page.route.id === '/[orgName]/[projectName]/runs/[runNumber]'}
			>
				<House data-slot="icon" />
				<SidebarLabel>Summary</SidebarLabel>
			</SidebarItem>

			<SidebarDivider />

			<SidebarSection>
				<SidebarHeading>Jobs</SidebarHeading>

				{#each sortedJobs as job (job.id)}
					<SidebarItem
						href={`${baseHref}/jobs/${job.number}`}
						current={page.route.id ===
							'/[orgName]/[projectName]/runs/[runNumber]/jobs/[jobNumber]' &&
							page.params.jobNumber === job.number.toString()}
					>
						<RunStatus class="size-5" status={job.status} conclusion={job.conclusion} />
						<SidebarLabel>
							{job.name}
						</SidebarLabel>
						<span class="text-on-surface-variant ml-auto whitespace-nowrap text-xs">
							<LiveDate finishedAt={job.finishedAt} startedAt={job.createdAt} />
						</span>
					</SidebarItem>
				{/each}
				{#if run.data && !run.data.jobs?.length}
					<Text class="pl-2">{run.data.conclusion ? 'No jobs' : 'Waiting for jobs'}</Text>
				{/if}

				{#if run.isLoading}
					<Skeleton class="my-1 h-7" />
					<Skeleton class="my-1 h-7" />
				{/if}
			</SidebarSection>

			<SidebarDivider />

			<SidebarSection>
				<SidebarHeading>Run details</SidebarHeading>

				<SidebarItem
					href={`${baseHref}/workflow`}
					current={page.route.id === '/[orgName]/[projectName]/runs/[runNumber]/workflow'}
				>
					<Notepad data-slot="icon" />
					<SidebarLabel>Workflow logs</SidebarLabel>
				</SidebarItem>
			</SidebarSection>
		</SidebarBody>
	</Sidebar>

	{@render children()}
</div>
