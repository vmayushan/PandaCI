<script lang="ts">
	import {
		SidebarItem,
		SidebarSection,
		ProjectSwitcher,
		SidebarLabel,
		SidebarDivider,
		SidebarLayout,
		SidebarHeading
	} from '$lib/components';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import type { Snippet } from 'svelte';
	import type { APIError } from '$lib/api';
	import { page } from '$app/state';
	import NotFound from '../notFound.svelte';
	import { GearSix, GitCommit, Lock, Stack } from 'phosphor-svelte';
	import ContainerLayout from '$lib/components/containerLayout.svelte';
	import clsx from 'clsx';
	import posthog from 'posthog-js';

	const { children }: { children: Snippet } = $props();

	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	const project = createQuery(() =>
		queries.organization.getByName(page.params.orgName)._ctx.projectByName(page.params.projectName)
	);

	$effect(() => {
		if (org.data && project.data) {
			posthog.resetGroups();
			posthog.group('org', org.data.id, {
				name: org.data.name,
				slug: org.data.slug
			});
			posthog.group('project', project.data.id, {
				name: project.data.name,
				slug: project.data.slug
			});
		}
	});
</script>

{#if project.isError && (project.error as APIError).status === 404}
	<NotFound itemName="Project" />
{:else}
	<SidebarLayout>
		{#snippet sidebarBody(collapsed, mobile)}
			<ProjectSwitcher disablePortal={mobile} />
			<SidebarDivider />
			<SidebarHeading class="pb-0.5">Project</SidebarHeading>
			<SidebarSection class={clsx('transform transition', collapsed && '-translate-y-6')}>
				<SidebarItem
					tooltip="Runs"
					current={page.route.id === '/[orgName]/[projectName]' ||
						page.route.id?.startsWith('/[orgName]/[projectName]/runs')}
					href={`/${page.params.orgName}/${page.params.projectName}`}
				>
					<GitCommit data-slot="icon" />
					<SidebarLabel>Runs</SidebarLabel>
				</SidebarItem>

				<SidebarItem
					tooltip="Variables"
					current={page.route.id === '/[orgName]/[projectName]/variables'}
					href={`/${page.params.orgName}/${page.params.projectName}/variables`}
				>
					<Lock data-slot="icon" />
					<SidebarLabel>Variables</SidebarLabel>
				</SidebarItem>

				<SidebarItem
					tooltip="Environments"
					current={page.route.id === '/[orgName]/[projectName]/environments'}
					href={`/${page.params.orgName}/${page.params.projectName}/environments`}
				>
					<Stack data-slot="icon" />
					<SidebarLabel>Environments</SidebarLabel>
				</SidebarItem>

				<SidebarItem
					tooltip="Settings"
					current={page.route.id === '/[orgName]/[projectName]/settings'}
					href={`/${page.params.orgName}/${page.params.projectName}/settings`}
				>
					<GearSix data-slot="icon" />
					<SidebarLabel>Settings</SidebarLabel>
				</SidebarItem>
			</SidebarSection>
		{/snippet}

		{#if page.route.id?.startsWith('/[orgName]/[projectName]/runs/[runNumber]')}
			{@render children()}
		{:else}
			<ContainerLayout>
				{@render children()}
			</ContainerLayout>
		{/if}
	</SidebarLayout>
{/if}
