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
	import CreateVariableModal from './createVariableModal.svelte';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import Tooltip from '$lib/components/tooltip/tooltip.svelte';
	import BadgeButton from '$lib/components/badgeButton.svelte';
	import VariableValue from './variableValue.svelte';

	dayjs.extend(relativeDate);

	const variables = createQuery(() =>
		queries.variables.projectVariables(page.params.orgName, page.params.projectName)._ctx.list()
	);

	const queryClient = useQueryClient();

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	const environments = createQuery(() =>
		queries.environments
			.projectEnvironments(page.params.orgName, page.params.projectName)
			._ctx.list()
	);

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	const deleteProjectVariableMutation = createMutation(() => ({
		mutationFn: (id: string) =>
			API.delete('/v1/orgs/{orgName}/projects/{projectName}/variables/{variableID}', {
				params: {
					orgName: page.params.orgName,
					projectName: page.params.projectName,
					variableID: id
				}
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.variables.projectVariables(page.params.orgName, page.params.projectName)._ctx.list()
			);
		}
	}));

	let createModalOpen = $state(false);
</script>

<Title title="Variables">
	{#snippet description()}
		Variables are loaded into workflow runs as environment variables.
	{/snippet}
	<Button onclick={() => (createModalOpen = true)}>New Variable</Button>
</Title>

{#if org.data && project.data && environments.data}
	<CreateVariableModal
		environments={environments.data}
		org={org.data}
		project={project.data}
		bind:open={createModalOpen}
	/>
{/if}

{#if variables.data?.length === 0}
	<Text class="mt-12 text-center">No variables have been created for this project yet.</Text>
{/if}

{#if variables.data?.length !== 0}
	<div class="mt-12 rounded-lg border border-zinc-950/10 px-6 py-8 sm:px-8 dark:border-white/10">
		<Table class="[--gutter:--spacing(6)] sm:[--gutter:--spacing(8)]">
			<TableHead>
				<TableRow>
					<TableHeader>Name</TableHeader>
					<TableHeader>Value</TableHeader>
					<TableHeader>Environments</TableHeader>
					<TableHeader>Created</TableHeader>
					<TableHeader class="relative w-0">
						<span class="sr-only">Actions</span>
					</TableHeader>
				</TableRow>
			</TableHead>
			<TableBody>
				{#if variables.data}
					{#each variables.data as variable (variable.id)}
						<TableRow>
							<TableCell class="font-medium">{variable.key}</TableCell>
							<TableCell>
								<VariableValue sensitive={variable.sensitive} id={variable.id} />
							</TableCell>
							<TableCell>
								{#if !variable.environments?.length}
									<Text class="text-zinc-500">All</Text>
								{:else}
									<ul class="flex space-x-2">
										{#each variable.environments as environment (environment.id)}
											<li>
												<BadgeButton
													color="zinc"
													href={`/${page.params.orgName}/${page.params.projectName}/environments`}
												>
													{environment.name}
												</BadgeButton>
											</li>
										{/each}
									</ul>
								{/if}
							</TableCell>
							<TableCell class="text-zinc-500">
								<Tooltip position="right" text={variable.createdAt}>
									{dayjs(variable.createdAt).fromNow()}
								</Tooltip>
							</TableCell>
							<TableCell>
								<Button
									outline
									onclick={() => {
										deleteProjectVariableMutation.mutate(variable.id);
									}}
								>
									Delete
								</Button>
							</TableCell>
						</TableRow>
					{/each}
				{/if}

				{#if variables.isLoading}
					{#each Array.from({ length: 10 }) as _, i (i)}
						<TableRow>
							{#each Array.from({ length: 5 }) as _, i (i)}
								<TableCell>
									<Skeleton class={clsx('h-9', i === 4 ? 'w-20' : 'w-full')} />
								</TableCell>
							{/each}
						</TableRow>
					{/each}
				{/if}
			</TableBody>
		</Table>
	</div>
{/if}
