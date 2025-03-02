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
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import { API } from '$lib/api';
	import clsx from 'clsx';
	import CreateEnvironmentModal from './createEnvironmentModal.svelte';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import Tooltip from '$lib/components/tooltip/tooltip.svelte';

	dayjs.extend(relativeDate);

	const environments = createQuery(() =>
		queries.environments
			.projectEnvironments(page.params.orgName, page.params.projectName)
			._ctx.list()
	);

	const queryClient = useQueryClient();

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	const deleteProjectEnvironmentsMutation = createMutation(() => ({
		mutationFn: (id: string) =>
			API.delete('/v1/orgs/{orgSlug}/projects/{projectSlug}/environments/{environmentID}', {
				params: {
					orgSlug: page.params.orgName,
					projectSlug: page.params.projectName,
					environmentID: id
				}
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.environments
					.projectEnvironments(page.params.orgName, page.params.projectName)
					._ctx.list()
			);
		}
	}));

	let createModalOpen = $state(false);
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
									outline
									onclick={() => {
										deleteProjectEnvironmentsMutation.mutate(environment.id);
									}}
								>
									Delete
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
