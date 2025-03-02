<script lang="ts">
	import { page } from '$app/state';
	import { API } from '$lib/api';
	import {
		Title,
		Button,
		Table,
		Skeleton,
		TableBody,
		TableCell,
		TableHeader,
		TableRow
	} from '$lib/components';
	import TableHead from '$lib/components/table/tableHead.svelte';
	import { queries } from '$lib/queries';
	import { getUser } from '$lib/runes/index.svelte';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import clsx from 'clsx';
	import InviteMemberModal from './inviteMemberModal.svelte';

	const queryClient = useQueryClient();

	const users = createQuery(() => queries.organization.getByName(page.params.orgName)._ctx.users());

	const deleteOrgUserMutation = createMutation(() => ({
		mutationFn: (id: string) =>
			API.delete('/v1/orgs/{orgName}/users/{userID}', {
				params: {
					orgName: page.params.orgName,
					userID: id
				}
			}),
		onSettled: () => {
			queryClient.invalidateQueries(
				queries.organization.getByName(page.params.orgName)._ctx.users()
			);
		}
	}));

	const currentUser = getUser();

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	let open = $state(false);
</script>

<Title title="Members">
	{#snippet description()}
		Invite users to your organization and manage their roles
	{/snippet}

	<Button onclick={() => (open = true)}>Invite member</Button>
</Title>

<InviteMemberModal bind:open />

<div class="mt-12 rounded-lg border border-zinc-950/10 px-6 py-8 sm:px-8 dark:border-white/10">
	<Table>
		<TableHead>
			<TableRow>
				<TableHeader>Name</TableHeader>
				<TableHeader>Email</TableHeader>
				<TableHeader>Role</TableHeader>
				<TableHeader class="relative w-0">
					<span class="sr-only">Actions</span>
				</TableHeader>
			</TableRow>
		</TableHead>
		<TableBody>
			{#each users.data ?? [] as user (user.id)}
				<TableRow>
					<TableCell class="font-medium">{user.name}</TableCell>
					<TableCell class="font-medium">{user.email}</TableCell>
					<TableCell class="font-medium">
						{user.role}
						{org.data?.ownerID === user.id ? '(Owner)' : ''}
					</TableCell>
					<TableCell>
						<Button
							disabled={user.id === org.data?.ownerID}
							plain
							onclick={() => deleteOrgUserMutation.mutate(user.id)}
						>
							{currentUser.data?.id === user.id ? 'Leave' : 'Remove'}
						</Button>
					</TableCell>
				</TableRow>
			{/each}
			{#if users.isLoading}
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
