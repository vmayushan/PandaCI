<script lang="ts">
	import { SidebarItem, SidebarLabel, SidebarSection } from '$lib/components/sidebar';
	import { BuildingOffice, User } from 'phosphor-svelte';
	import { page } from '$app/state';
	import SidebarLayout from '$lib/components/sidebarLayout.svelte';
	import ContainerLayout from '$lib/components/containerLayout.svelte';
	import SidebarHeading from '$lib/components/sidebar/sidebarHeading.svelte';
	import clsx from 'clsx';
	import { getUser } from '$lib/runes/user.svelte';
	import posthog from 'posthog-js';

	const { children } = $props();

	const user = getUser();

	$effect(() => {
		if (user.data) {
			posthog.identify(user.data.id, {
				email: user.data.email,
				name: user.data.name
			});
			posthog.resetGroups();
		}
	});
</script>

<SidebarLayout>
	{#snippet sidebarBody(collapsed)}
		<SidebarHeading class="pb-0.5">Account</SidebarHeading>
		<SidebarSection class={clsx('transform transition', collapsed && '-translate-y-6')}>
			<SidebarItem current={page.route.id === '/account/orgs'} href="/account">
				<BuildingOffice data-slot="icon" />
				<SidebarLabel>Organizations</SidebarLabel>
			</SidebarItem>
			<SidebarItem current={page.route.id === '/account/profile'} href="/account/profile">
				<User data-slot="icon" />
				<SidebarLabel>Profile</SidebarLabel>
			</SidebarItem>
		</SidebarSection>
	{/snippet}

	<ContainerLayout>
		{@render children()}
	</ContainerLayout>
</SidebarLayout>
