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
	<div class="mt-12 flex h-full flex-col text-center">
		<div class="my-auto">
			<Heading class="mb-2" level={3}>No Organizations</Heading>
			<Text>
				Get started by creating a new organization or asking an admin to invite you to one.
			</Text>

			<Button color="dark/white" class="mt-6" onclick={() => (open = true)}>
				Create Organization
			</Button>
		</div>
	</div>
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
			<Button loading={createOrgMutation.isPending} type="submit">Create organization</Button>
		</DialogActions>
	</form>
</Dialog>
