<script lang="ts">
	import { page } from '$app/state';
	import { queries } from '$lib/queries';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { CaretDown, GlobeHemisphereWest, Plus } from 'phosphor-svelte';
	import { Avatar, DropdownItem, Input, SidebarItem } from '..';
	import Dropdown from '../dropdown/dropdown.svelte';
	import DropdownButton from '../dropdown/dropdownButton.svelte';
	import DropdownMenu from '../dropdown/dropdownMenu.svelte';
	import { Text } from '../text';
	import DropdownSection from '../dropdown/dropdownSection.svelte';
	import Divider from '../divider.svelte';
	import { getContext } from 'svelte';
	import type { SidebarContext } from '../sidebar/sidebar.svelte';
	import Skeleton from '../skeleton.svelte';

	const currentOrg = createQuery(() => queries.organization.getByName(page.params.orgName!));

	const currentProject = createQuery(() => ({
		...queries.organization
			.getByName(page.params.orgName!)
			._ctx.projectByName(page.params.projectName),
		enabled: Boolean(page.params.projectName)
	}));

	const projects = createQuery(() => ({
		...queries.organization.getByName(page.params.orgName)._ctx.projects(),
		enabled: Boolean(currentOrg)
	}));

	const queryClient = useQueryClient();

	$effect(() => {
		if (projects.data) {
			projects.data.forEach((project) => {
				const key = queries.organization
					.getByName(page.params.orgName)
					._ctx.projectByName(project.slug).queryKey;
				if (!queryClient.getQueryData(key)) queryClient.setQueryData(key, project);
			});
		}
	});

	let search = $state('');

	const filteredProjects = $derived(
		projects.data?.filter(({ name }) =>
			name.toLowerCase().replaceAll(' ', '').includes(search.toLowerCase().replaceAll(' ', ''))
		)
	);

	const sidebarContext = getContext<SidebarContext>('sidebar');

	const { disablePortal }: { disablePortal: boolean } = $props();
</script>

<Dropdown>
	{@const selectedName = page.params.projectName
		? currentProject.data?.name
		: currentOrg.data?.name}
	<DropdownButton class={['-my-2', sidebarContext.collapsed && 'w-fit!']}>
		{#snippet renderer(props, iconProps)}
			<SidebarItem tooltip={sidebarContext.collapsed ? selectedName : undefined} {...props}>
				{#if page.params.projectName}
					<Avatar name={selectedName} src={currentProject.data?.avatarURL} alt="project picture" />
				{:else}
					<GlobeHemisphereWest data-slot="avatar" weight="duotone" />
				{/if}
				{#if !sidebarContext.collapsed}
					{#if selectedName}
						<span class="truncate">
							{selectedName}
						</span>
					{:else}
						<Skeleton class="h-4 w-20" />
					{/if}
				{/if}
				{#if !sidebarContext.collapsed}
					<CaretDown data-slot="icon" class="h-4 w-4 " {...iconProps} />
				{/if}
			</SidebarItem>
		{/snippet}
	</DropdownButton>
	<DropdownMenu class="max-w-sm overflow-hidden" {disablePortal}>
		<DropdownSection>
			<DropdownItem value="org-home" href={`/${page.params.orgName}`}>
				<GlobeHemisphereWest weight="duotone" data-slot="icon" />
				{currentOrg.data?.name}
			</DropdownItem>
		</DropdownSection>

		<Divider class="mb-2 mt-2" />

		<Input
			class="after:hidden! col-span-full my-0.5 items-center"
			autofocus
			bind:value={search}
			onkeydown={(e) => {
				if (e.key === ' ') e.stopPropagation();
			}}
			placeholder="Find project..."
			aria-label="Search projects"
		/>
		<DropdownSection class="mt-2 max-h-56 overflow-y-auto overflow-x-hidden">
			{#if filteredProjects?.length}
				{#each filteredProjects as project (project.id)}
					<DropdownItem value={project.id} href={`/${page.params.orgName}/${project.slug}`}>
						<Avatar name={project.name} alt="project icon" src={project.avatarURL} />
						<span class="w-full grow truncate">
							{project.name}
						</span>
					</DropdownItem>
				{/each}
			{:else}
				<Text class="mx-3 my-1.5 text-center">No projects found</Text>
			{/if}
		</DropdownSection>
		<Divider class="my-2" />
		<DropdownSection>
			<DropdownItem value="new-project" href={`/${page.params.orgName}/new`}>
				<Plus data-slot="icon" />New Project
			</DropdownItem>
		</DropdownSection>
	</DropdownMenu>
</Dropdown>
