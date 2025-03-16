<script lang="ts">
	import { page } from '$app/state';
	import {
		Title,
		Text,
		Dropdown,
		DropdownButton,
		DropdownItem,
		DropdownLabel,
		DropdownMenu,
		SidebarBody,
		SidebarItem,
		SidebarLabel,
		SidebarDivider,
		SidebarSection,
		SidebarHeading
	} from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { CaretLeft, House, Notepad, CaretDown } from 'phosphor-svelte';
	import RunStatus from '$lib/components/runStatus.svelte';
	import type { WorkflowRun } from '$lib/api/organization';
	import type { Snippet } from 'svelte';
	import Sidebar from '$lib/components/sidebar/sidebar.svelte';
	import Skeleton from '$lib/components/skeleton.svelte';
	import LiveDate from './liveDate.svelte';
	import DropdownSection from '$lib/components/dropdown/dropdownSection.svelte';
	import DropdownDivider from '$lib/components/dropdown/dropdownDivider.svelte';
	import DropdownHeading from '$lib/components/dropdown/dropdownHeading.svelte';

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

	const currentNav = $derived.by(() => {
		if (page.route.id === '/[orgName]/[projectName]/runs/[runNumber]/workflow')
			return {
				label: 'Workflow logs',
				icon: Notepad
			};

		if (page.route.id === '/[orgName]/[projectName]/runs/[runNumber]/jobs/[jobNumber]') {
			const job = run.data?.jobs?.find((job) => job.number.toString() === page.params.jobNumber);
			return {
				label: job?.name ?? 'Unknown job',
				conclusion: job?.conclusion,
				status: job?.status
			};
		}

		return {
			label: 'Summary',
			icon: House
		};
	});
</script>

<Title
	titleLoading={run.isLoading}
	state={{
		status: run.data?.status,
		conclusion: run.data?.conclusion
	}}
	title={run.data?.gitTitle ?? `#${run.data?.number}`}
>
	{#snippet description()}
		{#if run.isLoading}
			<Skeleton class="ml-10 mt-3.5 h-3.5 w-56" />
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

<div class="mt-4 flex h-fit flex-col space-x-4 lg:flex-row">
	<Dropdown>
		<DropdownButton outline class="mb-4 w-min whitespace-nowrap lg:hidden">
			{#snippet indicator(props)}
				<CaretDown data-slot="icon" {...props} />
			{/snippet}
			{#if currentNav?.icon}
				<currentNav.icon class="size-4" data-slot="icon" />
			{/if}
			{#if currentNav.status}
				<RunStatus class="size-4" status={currentNav.status} conclusion={currentNav.conclusion} />
			{/if}
			{currentNav?.label}
		</DropdownButton>
		<DropdownMenu>
			<DropdownSection>
				<DropdownItem
					href={`/${page.params.orgName}/${page.params.projectName}/runs/${page.params.runNumber}`}
					value="summary"
				>
					<House data-slot="icon" />
					<DropdownLabel>Summary</DropdownLabel>
				</DropdownItem>
			</DropdownSection>

			<DropdownDivider />

			<DropdownSection>
				<DropdownHeading>Jobs</DropdownHeading>
				{#each sortedJobs as job (job.id)}
					<DropdownItem href={`${baseHref}/jobs/${job.number}`} value={`job-${job.id}`}>
						<RunStatus
							data-slot="icon"
							class="size-5"
							status={job.status}
							conclusion={job.conclusion}
						/>
						<DropdownLabel>
							{job.name}
						</DropdownLabel>

						<span
							class="text-on-surface-variant col-start-5 row-start-1 flex justify-self-end whitespace-nowrap text-xs"
						>
							<LiveDate finishedAt={job.finishedAt} startedAt={job.createdAt} />
						</span>
					</DropdownItem>
				{/each}
				{#if run.data && !run.data.jobs?.length}
					<Text class="pl-3">{run.data.conclusion ? 'No jobs' : 'Waiting for jobs'}</Text>
				{/if}

				{#if run.isLoading}
					<div class="col-span-4 px-3.5">
						<Skeleton class="my-1 h-7 w-full" />
						<Skeleton class="my-1 h-7 w-full" />
					</div>
				{/if}
			</DropdownSection>

			<DropdownDivider />

			<DropdownSection>
				<DropdownItem
					href={`/${page.params.orgName}/${page.params.projectName}/runs/${page.params.runNumber}/workflow`}
					value="workflow-logs"
				>
					<Notepad data-slot="icon" />
					<DropdownLabel>Workflow logs</DropdownLabel>
				</DropdownItem>
			</DropdownSection>
		</DropdownMenu>
	</Dropdown>

	<Sidebar class="sticky top-0 hidden w-full max-w-56 shrink-0 grow flex-col lg:flex 2xl:max-w-xs">
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
