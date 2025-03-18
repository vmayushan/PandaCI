<script lang="ts">
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import type { GitInstallation, GitRepository } from '$lib/api';
	import { Button, Description, Subheading, Text, Heading, Skeleton } from '$lib/components';
	import { Lock } from 'phosphor-svelte';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import CreateProjectModal from './createProjectModal.svelte';
	import { PUBLIC_GITHUB_APP_NAME } from '$env/static/public';

	dayjs.extend(relativeTime);

	interface Props {
		installation: GitInstallation;
		search: string;
	}

	let { installation = $bindable(), search = $bindable() }: Props = $props();

	let debouncedSearch = $state('');

	let timer: ReturnType<typeof setTimeout> | undefined;

	$effect(() => {
		clearTimeout(timer);
		// eslint-disable-next-line @typescript-eslint/no-unused-expressions
		search;

		timer = setTimeout(() => {
			debouncedSearch = search;
		}, 300);
	});

	const allRepos = createQuery(() => ({
		...queries.github.listInstallation()._ctx.listRepositories(installation.id),
		enabled: installation.repositoryScopes === 'selected'
	}));

	const query = $derived(
		`${installation.isUser ? 'user' : 'org'}:"${installation.name}""${debouncedSearch}"`
	);

	const searchedRepos = createQuery(() => ({
		...queries.github.listInstallation()._ctx.listRepositories(installation.id, { query }),
		enabled: installation.repositoryScopes === 'all'
	}));

	const namedRepoEnabled = $derived(
		installation.repositoryScopes === 'selected' &&
			Boolean(allRepos.data?.limitExceeded) &&
			Boolean(debouncedSearch)
	);

	const namedRepo = createQuery(() => ({
		...queries.github
			.listInstallation()
			._ctx.listRepositories(installation.id, { name: debouncedSearch, owner: installation.name }),
		enabled: namedRepoEnabled
	}));

	const repos = $derived(
		(namedRepo.data?.repos || allRepos.data?.repos || searchedRepos.data?.repos || [])
			.filter(
				(repo) =>
					installation.repositoryScopes === 'all' ||
					repo.name
						.toLowerCase()
						.replaceAll(' ', '')
						.includes(debouncedSearch.toLowerCase().replaceAll(' ', ''))
			)
			.slice(0, 10)
	);

	const isLoading = $derived(allRepos.isLoading || searchedRepos.isLoading || namedRepo.isLoading);

	let open = $state(false);

	let repository = $state<GitRepository>();
</script>

{#if repository}
	<CreateProjectModal bind:open repo={repository} {installation} />
{/if}

{#if repos.length === 0 && !allRepos.isLoading && !searchedRepos.isLoading && !namedRepo.isLoading}
	<div class="flex flex-col justify-center">
		<Heading class="text-center">No repos found</Heading>
		<Text class="text-center">Make sure you've allowed our app to access your repositories</Text>
		<Button
			color="dark/white"
			href={`https://github.com/apps/${PUBLIC_GITHUB_APP_NAME}/installations/new`}
			class="mx-auto mt-4"
		>
			Update permissions
		</Button>
	</div>
{/if}

{#if !isLoading}
	<ul role="list" class="divide-y divide-zinc-950/5 dark:divide-white/5">
		{#each repos as repo (repo.id)}
			<li class="flex items-center justify-between gap-x-6 py-5">
				<div class="min-w-0">
					<div class="flex items-center gap-x-3">
						<Subheading>{repo.name}</Subheading>
						{#if !repo.public}<Lock data-slot="icon" class="dark:text-zinc-40 text-zinc-500" />{/if}
					</div>
					<div class="mt-1 flex items-center gap-x-2 text-xs leading-5 text-zinc-500">
						<Text class="whitespace-nowrap">
							Updated <time class="whitespace-nowrap" datetime={repo.updatedAt}>
								{dayjs().to(dayjs(repo.updatedAt))}
							</time>
						</Text>
						{#if repo.description}
							<svg viewBox="0 0 2 2" class="h-0.5 w-0.5 shrink-0 fill-current">
								<circle cx="1" cy="1" r="1" />
							</svg>
							<Description class="truncate">{repo.description}</Description>
						{/if}
					</div>
				</div>
				<div class="flex flex-none items-center gap-x-4">
					<Button
						color="dark/white"
						onclick={() => {
							repository = repo;
							open = true;
						}}
					>
						Import
					</Button>
				</div>
			</li>
		{/each}
	</ul>
{/if}

{#if isLoading}
	<div class="flex flex-col justify-center">
		{#each Array.from({ length: 5 }) as _, i (i)}
			<Skeleton class="my-4 h-16 w-full" />
		{/each}
	</div>
{/if}
