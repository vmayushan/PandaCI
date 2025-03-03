<script lang="ts">
	import {
		Title,
		Text,
		Button,
		Table,
		TableRow,
		TableCell,
		Skeleton,
		TableHeader,
		TableHead,
		TableBody
	} from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import clsx from 'clsx';
	import CreateEnvironmentModal from './createEnvironmentModal.svelte';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import Tooltip from '$lib/components/tooltip/tooltip.svelte';
	import EditEnvironmentModal from './editEnvironmentModal.svelte';
	import { PencilSimple, TrashSimple } from 'phosphor-svelte';
	import type { ProjectEnvironment } from '$lib/api/organization';
	import DeleteEnvironmentModal from './deleteEnvironmentModal.svelte';

	dayjs.extend(relativeDate);

	const environments = createQuery(() =>
		queries.environments
			.projectEnvironments(page.params.orgName, page.params.projectName)
			._ctx.list()
	);

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	let createModalOpen = $state(false);

	let editEnvironment = $state<ProjectEnvironment | undefined>();
	let deleteEnvironment = $state<ProjectEnvironment | undefined>();
</script>

<Title title="Environments">
	<Button onclick={() => (createModalOpen = true)}>New Environment</Button>
	{#snippet description()}
		Environments are used to restrict variables to specific branches.
	{/snippet}
</Title>

{#if org.data && project.data}
	<CreateEnvironmentModal org={org.data} project={project.data} bind:open={createModalOpen} />
{/if}

{#if org.data && project.data}
	<EditEnvironmentModal
		org={org.data}
		project={project.data}
		bind:oldEnvironment={editEnvironment}
	/>
{/if}

{#if org.data && project.data}
	<DeleteEnvironmentModal
		org={org.data}
		project={project.data}
		bind:environment={deleteEnvironment}
	/>
{/if}

{#if environments.data?.length === 0}
	<Text class="mt-12 text-center">No environments have been created for this project yet.</Text>
{/if}

{#if environments.data?.length !== 0}
	<div class="mt-12 rounded-lg border border-zinc-950/10 px-6 py-8 sm:px-8 dark:border-white/10">
		<Table class="[--gutter:--spacing(6)] sm:[--gutter:--spacing(8)]">
			<TableHead>
				<TableRow>
					<TableHeader>Name</TableHeader>
					<TableHeader>Branch Pattern</TableHeader>
					<TableHeader>Created</TableHeader>
					<TableHeader class="relative w-0">
						<span class="sr-only">Actions</span>
					</TableHeader>
				</TableRow>
			</TableHead>
			<TableBody>
				{#if environments.data}
					{#each environments.data as environment (environment.id)}
						<TableRow>
							<TableCell class="font-medium">{environment.name}</TableCell>
							<TableCell>{environment.branchPattern}</TableCell>
							<TableCell class="text-zinc-500">
								<Tooltip position="right" text={environment.createdAt}>
									{dayjs(environment.createdAt).fromNow()}
								</Tooltip>
							</TableCell>
							<TableCell>
								<Button
									plain
									tooltip="Delete"
									onclick={() => {
										deleteEnvironment = environment;
									}}
								>
									<span class="sr-only">Delete {environment.name}</span>
									<TrashSimple weight="fill" data-slot="icon" />
								</Button>

								<Button
									plain
									tooltip="Edit"
									onclick={() => {
										editEnvironment = environment;
									}}
								>
									<span class="sr-only">Edit {environment.name}</span>
									<PencilSimple weight="fill" data-slot="icon" />
								</Button>
							</TableCell>
						</TableRow>
					{/each}
				{/if}

				{#if environments.isLoading}
					{#each Array.from({ length: 10 }) as _, i (i)}
						<TableRow>
							{#each Array.from({ length: 4 }) as _, i (i)}
								<TableCell>
									<Skeleton class={clsx('h-9', i === 3 ? 'w-20' : 'w-full')} />
								</TableCell>
							{/each}
						</TableRow>
					{/each}
				{/if}
			</TableBody>
		</Table>
	</div>
{/if}
