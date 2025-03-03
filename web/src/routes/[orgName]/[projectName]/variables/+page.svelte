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
	import CreateVariableModal from './createVariableModal.svelte';
	import dayjs from 'dayjs';
	import relativeDate from 'dayjs/plugin/relativeTime';
	import Tooltip from '$lib/components/tooltip/tooltip.svelte';
	import BadgeButton from '$lib/components/badgeButton.svelte';
	import VariableValue from './variableValue.svelte';
	import EditVariableModal from './editVariableModal.svelte';
	import { type ProjectVariable } from '$lib/api/organization';
	import { PencilSimple, TrashSimple } from 'phosphor-svelte';
	import DeleteVariableModal from './deleteVariableModal.svelte';

	dayjs.extend(relativeDate);

	const variables = createQuery(() =>
		queries.variables.projectVariables(page.params.orgName, page.params.projectName)._ctx.list()
	);

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	const environments = createQuery(() =>
		queries.environments
			.projectEnvironments(page.params.orgName, page.params.projectName)
			._ctx.list()
	);

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	let createModalOpen = $state(false);

	let editVariable = $state<ProjectVariable | undefined>();

	let deleteVariable = $state<ProjectVariable | undefined>();
</script>

<Title title="Variables">
	{#snippet description()}
		Variables are loaded into workflow runs as environment variables.
	{/snippet}
	<Button onclick={() => (createModalOpen = true)}>New Variable</Button>
</Title>

{#if org.data && project.data}
	<CreateVariableModal
		environments={environments.data ?? []}
		org={org.data}
		project={project.data}
		bind:open={createModalOpen}
	/>
{/if}

{#if org.data && project.data}
	<DeleteVariableModal org={org.data} project={project.data} bind:variable={deleteVariable} />
{/if}

{#if org.data && project.data}
	<EditVariableModal
		environments={environments.data ?? []}
		org={org.data}
		project={project.data}
		bind:oldVariable={editVariable}
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
									plain
									tooltip="Delete"
									onclick={() => {
										deleteVariable = variable;
									}}
								>
									<span class="sr-only">Delete {variable.key}</span>
									<TrashSimple weight="fill" data-slot="icon" />
								</Button>

								<Button
									plain
									tooltip="Edit"
									onclick={() => {
										editVariable = variable;
									}}
								>
									<span class="sr-only">Edit {variable.key}</span>
									<PencilSimple weight="fill" data-slot="icon" />
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
