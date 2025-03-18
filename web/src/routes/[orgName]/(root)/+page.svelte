<script lang="ts">
	import { page } from '$app/state';
	import { Button, Title, Skeleton } from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import { FolderPlus } from 'phosphor-svelte';

	dayjs.extend(relativeDate);

	const projects = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projects()
	);

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));
</script>

<Title title="Projects">
	<Button href={`/${page.params.orgName}/new`}>New project</Button>
	{#snippet description()}
		Projects contain your workflow runs, pulling their configurations from your repository.
	{/snippet}
</Title>

{#snippet emptyProjects()}
	<a
		href={`/${page.params.orgName}/new`}
		class="border-outline-variant hover:border-outline ring-offset-surface focus:outline-hidden relative mx-auto mt-12 block w-full max-w-2xl rounded-lg border-2 border-dashed p-12 text-center focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
	>
		<FolderPlus class="mx-auto size-12 text-zinc-400 dark:text-zinc-500" />
		<span class="mt-2 block text-sm font-semibold text-zinc-900 dark:text-white">
			Create a new project
		</span>
	</a>
{/snippet}

{#snippet projectList()}
	<ul role="list" class="divide-outline-variant divide-y">
		{#each projects.data ?? [] as project (project.id)}
			<li class="relative flex items-center space-x-4 py-4">
				<div class="min-w-0 flex-auto">
					<div class="flex items-center gap-x-3">
						<h2 class="min-w-0 text-sm font-semibold leading-6 text-zinc-950 dark:text-white">
							<a href={`/${org.data?.slug}/${project.slug}`} class="flex gap-x-2">
								<span class="truncate">{project.name}</span>
								<span class="absolute inset-0"></span>
							</a>
						</h2>
					</div>
					<div
						class="mt-3 flex items-center gap-x-2.5 text-xs leading-5 text-zinc-700 dark:text-zinc-400"
					>
						<p class="truncate">Deploys from GitHub</p>
						{#if project.lastBuild}
							<svg viewBox="0 0 2 2" class="h-0.5 w-0.5 flex-none fill-zinc-600 dark:fill-zinc-300">
								<circle cx="1" cy="1" r="1" />
							</svg>
							<p class="whitespace-nowrap">Last build {dayjs(project.lastBuild).fromNow()}</p>
						{/if}
					</div>
				</div>
				<svg
					class="h-5 w-5 flex-none text-zinc-400"
					viewBox="0 0 20 20"
					fill="currentColor"
					aria-hidden="true"
				>
					<path
						fill-rule="evenodd"
						d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
						clip-rule="evenodd"
					/>
				</svg>
			</li>
		{/each}
	</ul>
{/snippet}
{#if projects.isLoading}
	<div class="flex flex-col space-y-4 overflow-hidden pt-4">
		{#each Array.from({ length: 10 }) as _, i (i)}
			<Skeleton class="h-[4.5rem] w-full" />
		{/each}
	</div>
{:else if projects.data?.length}
	{@render projectList()}
{:else}
	{@render emptyProjects()}
{/if}
