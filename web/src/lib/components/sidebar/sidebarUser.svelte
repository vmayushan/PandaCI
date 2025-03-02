<script lang="ts">
	import { ArrowLineRight, BuildingOffice, CaretUp, User } from 'phosphor-svelte';
	import { getCurrentOrg, getUser } from '$lib/runes/index.svelte';
	import { SidebarItem } from '.';
	import { Avatar } from '..';
	import { Dropdown, DropdownButton, DropdownMenu, DropdownItem, DropdownLabel } from '../dropdown';
	import { page } from '$app/state';
	import { type SidebarContext } from './sidebar.svelte';
	import { getContext } from 'svelte';
	import clsx from 'clsx';

	const user = getUser();

	const currentOrg = getCurrentOrg({
		page
	});

	let windowWidth = $state(0);

	const sidebarContext = getContext<SidebarContext>('sidebar');
</script>

<svelte:window bind:innerWidth={windowWidth} />

<Dropdown>
	<DropdownButton>
		{#snippet renderer(props, iconProps)}
			{@const selectedName = page.params.orgName ? currentOrg.data?.name : 'My Account'}
			<SidebarItem
				class={clsx(
					sidebarContext.collapsed && 'w-fit -translate-x-2 py-0.5',
					'transform transition'
				)}
				tooltip={selectedName}
				{...props}
			>
				<span class="flex min-w-0 items-center gap-3 text-left">
					<Avatar
						name={page.params.orgName ? currentOrg.data?.name : user.data?.name}
						src={page.params.orgName ? currentOrg.data?.avatarURL : user.data?.avatar}
						class="size-9"
						square
						alt="profile"
					/>
					{#if !sidebarContext.collapsed}
						<span class="min-w-0">
							<span class="block truncate text-sm/5 font-medium text-zinc-950 dark:text-white">
								{page.params.orgName ? currentOrg.data?.name : 'My Account'}
							</span>
							<span class="block truncate text-xs/5 font-normal text-zinc-500 dark:text-zinc-400">
								{user.data?.name}
							</span>
						</span>
					{/if}
				</span>

				{#if !sidebarContext.collapsed}
					<CaretUp class="h-4 w-4 " {...iconProps} />
				{/if}
			</SidebarItem>
		{/snippet}
	</DropdownButton>
	<DropdownMenu disablePortal={windowWidth < 1024}>
		<DropdownItem href="/account/profile" value="profile">
			<User data-slot="icon" />
			<DropdownLabel>Profile</DropdownLabel>
		</DropdownItem>
		<DropdownItem href="/account/orgs" value="orgs">
			<BuildingOffice data-slot="icon" />
			<DropdownLabel>Switch organization</DropdownLabel>
		</DropdownItem>
		<DropdownItem href="/logout" value="logout">
			<ArrowLineRight data-slot="icon" />
			<DropdownLabel>Logout</DropdownLabel>
		</DropdownItem>
	</DropdownMenu>
</Dropdown>
