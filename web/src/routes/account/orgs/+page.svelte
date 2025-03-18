<script lang="ts">
	import {
		Button,
		Heading,
		Skeleton,
		Table,
		TableBody,
		TableCell,
		TableHead,
		TableHeader,
		TableRow,
		Text,
		Title,
		Dialog,
		DialogActions,
		DialogBody,
		DialogCloseButton,
		DialogDescription,
		DialogTitle
	} from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import CreateOrgForm, { type CreateOrgFormData } from '../../welcome/createOrgForm.svelte';
	import { handleForm } from '$lib/utils';
	import { useCreateOrgMutation } from '../useCreateOrg';
	import clsx from 'clsx';
	import { getUser } from '$lib/runes/user.svelte';
	import { BuildingOffice } from 'phosphor-svelte';

	const orgs = createQuery(() => queries.organization.list());

	const user = getUser();

	let open = $state(false);

	const createOrgMutation = useCreateOrgMutation();
</script>

<Title title="Organizations">
	<Button class="" onclick={() => (open = true)}>Create Organization</Button>
	{#snippet description()}
		Orgs you're a member of are listed here.
	{/snippet}
</Title>

{#if !orgs.data?.length && !orgs.isLoading}
	<button
		onclick={() => (open = true)}
		class="border-outline-variant hover:border-outline ring-offset-surface focus:outline-hidden relative mx-auto mt-12 block w-full max-w-2xl cursor-pointer rounded-lg border-2 border-dashed p-12 text-center focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
	>
		<BuildingOffice class="mx-auto size-12 text-zinc-400 dark:text-zinc-500" />
		<Heading level={2} class="mt-2 block text-sm font-semibold text-zinc-900 dark:text-white">
			Create a free org to get started.
		</Heading>
		<Text>Or you can ask an admin to invite you to one.</Text>
	</button>
{/if}

{#if orgs.data?.length || orgs.isLoading}
	<div class="mt-12 rounded-lg border border-zinc-950/10 px-6 py-8 sm:px-8 dark:border-white/10">
		<Table>
			<TableHead>
				<TableRow>
					<TableHeader>Name</TableHeader>
					<TableHeader>Role</TableHeader>
					<TableHeader>Plan</TableHeader>
				</TableRow>
			</TableHead>
			<TableBody>
				{#each orgs.data ?? [] as org (org.id)}
					<TableRow href={`/${org.slug}`}>
						<TableCell class="font-medium">{org.name}</TableCell>
						<TableCell>
							{org.currentUsersRole}
							{#if org.ownerID === user.data?.id}
								<span class="text-on-surface-variant">(owner)</span>
							{/if}
						</TableCell>
						<TableCell class="first-letter:uppercase">{org.license?.plan}</TableCell>
					</TableRow>
				{/each}
				{#if orgs.isLoading}
					{#each Array.from({ length: 3 }) as _, i (i)}
						<TableRow>
							{#each Array.from({ length: 2 })}
								<TableCell>
									<Skeleton class={clsx('h-9', 'w-full')} />
								</TableCell>
							{/each}
						</TableRow>
					{/each}
				{/if}
			</TableBody>
		</Table>
	</div>
{/if}

<Dialog bind:open>
	<DialogTitle>Create your organization</DialogTitle>
	<DialogDescription>Set up your organization's workspace on our platform.</DialogDescription>

	{#if createOrgMutation.error}
		<Text variant="error">{createOrgMutation.error.message}</Text>
	{/if}

	<form
		onsubmit={(e) => {
			const { data } = handleForm<CreateOrgFormData>(e);

			createOrgMutation.mutate(data);
		}}
	>
		<DialogBody>
			<CreateOrgForm />
		</DialogBody>

		<DialogActions>
			<DialogCloseButton plain>Cancel</DialogCloseButton>
			<Button color="dark/white" loading={createOrgMutation.isPending} type="submit">
				Create organization
			</Button>
		</DialogActions>
	</form>
</Dialog>
