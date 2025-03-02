<script lang="ts">
	import { page } from '$app/state';
	import {
		SidebarItem,
		SidebarSection,
		ProjectSwitcher,
		SidebarLabel,
		SidebarHeading
	} from '$lib/components';
	import ContainerLayout from '$lib/components/containerLayout.svelte';
	import SidebarDivider from '$lib/components/sidebar/SidebarDivider.svelte';
	import SidebarLayout from '$lib/components/sidebarLayout.svelte';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import clsx from 'clsx';
	import { CreditCard, GearSix, GitBranch, UserList } from 'phosphor-svelte';
	import posthog from 'posthog-js';
	import type { Snippet } from 'svelte';

	const { children }: { children: Snippet } = $props();

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	$effect(() => {
		if (org.data) {
			posthog.resetGroups();
			posthog.group('org', org.data.id, {
				name: org.data.name,
				slug: org.data.slug
			});
		}
	});
</script>

<SidebarLayout>
	{#snippet sidebarBody(collapsed, mobile)}
		<ProjectSwitcher disablePortal={mobile} />
		<SidebarDivider />
		<SidebarHeading class="pb-0.5">Organization</SidebarHeading>
		<SidebarSection class={clsx('transform transition', collapsed && '-translate-y-6')}>
			<SidebarItem
				tooltip="Projects"
				current={page.route.id === '/[orgName]/(root)'}
				href={`/${org.data?.slug}`}
			>
				<GitBranch data-slot="icon" />
				<SidebarLabel>Projects</SidebarLabel>
			</SidebarItem>

			<SidebarItem
				tooltip="Members"
				current={page.route.id === '/[orgName]/(root)/~/members'}
				href={`/${org.data?.slug}/~/members`}
			>
				<UserList data-slot="icon" />
				<SidebarLabel>Members</SidebarLabel>
			</SidebarItem>
			<SidebarItem
				tooltip="Billing"
				current={page.route.id === '/[orgName]/(root)/~/billing'}
				href={`/${org.data?.slug}/~/billing`}
			>
				<CreditCard data-slot="icon" />
				<SidebarLabel>Billing</SidebarLabel>
			</SidebarItem>

			<SidebarItem
				tooltip="Settings"
				current={page.route.id === '/[orgName]/(root)/~/settings'}
				href={`/${org.data?.slug}/~/settings`}
			>
				<GearSix data-slot="icon" />
				<SidebarLabel>Settings</SidebarLabel>
			</SidebarItem>
		</SidebarSection>
	{/snippet}

	<ContainerLayout>{@render children()}</ContainerLayout>
</SidebarLayout>
