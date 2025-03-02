<script lang="ts">
	import { page } from '$app/state';
	import { Button, TextLink, Title } from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery, createInfiniteQuery } from '@tanstack/svelte-query';
	import TriggerRunModal from './triggerRunModal.svelte';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import RunStatus from '$lib/components/runStatus.svelte';
	import Skeleton from '$lib/components/skeleton.svelte';
	import Text from '$lib/components/text/text.svelte';
	import { API } from '$lib/api';
	import Avatar from '$lib/components/avatar.svelte';
	import { GitBranch, GitCommit } from 'phosphor-svelte';
	import LiveDate from './runs/[runNumber]/liveDate.svelte';

	dayjs.extend(relativeDate);

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	const runs = createInfiniteQuery(() => ({
		queryKey: queries.runs.projectRuns(page.params.orgName, page.params.projectName).queryKey,
		queryFn: ({ pageParam }) =>
			API.get('/v1/orgs/{orgSlug}/projects/{projectSlug}/runs', {
				params: {
					orgSlug: page.params.orgName,
					projectSlug: page.params.projectName
				},
				queries: {
					page: pageParam
				}
			}),
		refetchInterval: 30_000,
		initialPageParam: 0,
		getNextPageParam: (lastPage, _, prevPageParam) => {
			if (lastPage.next) {
				return prevPageParam + 1;
			}
			return undefined;
		},
		getPreviousPageParam: (_, __, prevPageParam) => {
			if (prevPageParam > 0) {
				return prevPageParam - 1;
			}

			return undefined;
		}
	}));

	let openTriggerModal = $state(false);
</script>

<Title title="Runs">
	<Button onclick={() => (openTriggerModal = true)}>Trigger</Button>
	{#snippet description()}
		Runs are the executions of your workflows. See our <TextLink
			target="_blank"
			href="https://pandaci.com/docs/typescript-syntax/api/overview">docs</TextLink
		> for more information.
	{/snippet}
</Title>

{#if project.data && org.data}
	<TriggerRunModal project={project.data} org={org.data} bind:open={openTriggerModal} />
{/if}

{#if runs.isLoading}
	<div class="flex flex-col space-y-4 overflow-hidden pt-4">
		{#each Array.from({ length: 10 }) as _, i (i)}
			<Skeleton class="h-[4.5rem] w-full" />
		{/each}
	</div>
{/if}

{#if runs.data?.pages[0].data.length === 0}
	<Text class="mt-8 text-center">No runs have been triggered for this project yet.</Text>
{/if}

{#if runs.data}
	<ul role="list" class="divide-outline-variant divide-y">
		{#each runs.data.pages as runPage, i (i)}
			{#each runPage.data as run (run.id)}
				<li class="relative flex items-center space-x-4 py-4">
					<div class="min-w-0 flex-auto">
						<div class="flex items-center gap-x-3">
							<RunStatus class="size-5" status={run.status} conclusion={run.conclusion} />
							<h2 class="min-w-0 text-sm/6 font-bold text-zinc-900 dark:text-white">
								<a
									href={`/${page.params.orgName}/${page.params.projectName}/runs/${run.number}`}
									class="flex gap-x-2"
								>
									<span class="truncate">{run.gitTitle ?? '#' + run.number}</span>
									<span class="absolute inset-0"></span>
								</a>
							</h2>
							<span
								class="text-on-surface-variant ml-auto flex items-center justify-center space-x-2 text-xs"
							>
								{#if run.status === 'completed'}
									Ran for
								{/if}
								<LiveDate startedAt={run.createdAt} finishedAt={run.finishedAt} />
							</span>
						</div>
						<div
							class="mt-3 flex items-center gap-x-2.5 text-xs/5 text-zinc-700 dark:text-zinc-400"
						>
							<p class="font-semibold">
								{run.name}
							</p>
							<p>
								#{run.number}
							</p>

							<Avatar
								class="size-4"
								alt="commiter"
								name={run.committer?.name || run.committer?.email}
								src={run.committer?.avatar}
							/>
							{#if run.gitBranch}
								<GitBranch />
								<p class="truncate">{run.gitBranch}</p>
							{/if}
							<GitCommit class="" />
							<p class="truncate">{run.gitSha.slice(0, 7)}</p>
							<p class="ml-auto whitespace-nowrap">Initiated {dayjs(run.createdAt).fromNow()}</p>
						</div>
					</div>
					<svg
						class="h-5 w-5 flex-none text-zinc-500 dark:text-zinc-400"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
						data-slot="icon"
					>
						<path
							fill-rule="evenodd"
							d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z"
							clip-rule="evenodd"
						/>
					</svg>
				</li>
			{/each}
		{/each}
	</ul>
	{#if runs.hasNextPage}
		<div class="mt-8 flex w-full items-center justify-center">
			<Button plain onclick={() => runs.fetchNextPage()}>Load more</Button>
		</div>
	{/if}
{/if}
