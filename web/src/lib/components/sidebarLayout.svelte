<script lang="ts">
	import type { Snippet } from 'svelte';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { useMachine, normalizeProps, portal } from '@zag-js/svelte';
	import { nanoid } from 'nanoid';
	import * as dialog from '@zag-js/dialog';
	import { List, X, Newspaper, TextOutdent, TextIndent } from 'phosphor-svelte';
	import Button from './button.svelte';
	import { fade, fly } from 'svelte/transition';
	import { page } from '$app/state';
	import {
		SidebarHeader,
		SidebarBody,
		SidebarSection,
		SidebarItem,
		SidebarLabel,
		SidebarSpacer,
		SidebarHeading,
		SidebarFooter,
		SidebarUser,
		Sidebar
	} from './sidebar';
	import Verify from './verify.svelte';
	import Feedback from './feedback.svelte';
	import OrgPlanBanner from './orgPlanBanner.svelte';

	interface SidebarLayoutProps {
		children: Snippet;
		sidebarBody: Snippet<[boolean, boolean]>;
	}

	const { children, sidebarBody }: SidebarLayoutProps = $props();

	const id = nanoid(6);

	const service = useMachine(dialog.machine, { id });

	const api = $derived(dialog.connect(service, normalizeProps));

	const orgs = createQuery(() => queries.organization.list());

	let oldPathname = page.url.pathname;

	$effect(() => {
		if (oldPathname !== page.url.pathname) {
			api.setOpen(false);
			oldPathname = page.url.pathname;
		}
	});

	let collapsed = $state(localStorage.getItem('sidebarCollapsed') === 'true');

	$effect(() => {
		localStorage.setItem('sidebarCollapsed', collapsed.toString());
	});
</script>

{#snippet sidebar(mobile: boolean = false)}
	<Sidebar collapsed={mobile ? false : collapsed}>
		<SidebarHeader>
			<div class="flex items-center justify-center">
				{#if mobile || !collapsed}
					<a
						href={`/${page.params.orgName ?? (orgs.data?.length === 1 ? orgs.data?.[0]?.slug : 'account/orgs')}`}
						class="focus:outline-hidden ml-2 w-full rounded-lg text-2xl/8 font-semibold text-zinc-950 focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 sm:text-xl/8 dark:text-white dark:focus-visible:ring-offset-zinc-950"
					>
						Panda CI
					</a>
				{/if}
				{#if mobile}
					<Button plain aria-label="Close navigation" {...api.getCloseTriggerProps() as any}>
						<X data-slot="icon" />
					</Button>
				{:else}
					<Button
						class={[collapsed ? '-ml-1' : '-mr-0.5']}
						plain
						aria-label={collapsed ? 'Expand navigation' : 'Collapse navigation'}
						aria-expanded={!collapsed}
						onclick={() => (collapsed = !collapsed)}
					>
						{#if collapsed}
							<TextIndent data-slot="icon" />
						{:else}
							<TextOutdent data-slot="icon" />
						{/if}
					</Button>
				{/if}
			</div>
		</SidebarHeader>
		<SidebarBody>
			{@render sidebarBody(mobile ? false : collapsed, mobile)}
			<SidebarSpacer />
			<SidebarSection>
				<SidebarHeading>Resources</SidebarHeading>

				<Feedback {mobile} />
				<SidebarItem
					tooltip="Documentation"
					href="https://pandaci.com/docs/platform/overview/quick-start"
					target="_blank"
				>
					<Newspaper data-slot="icon" /><SidebarLabel>Docs</SidebarLabel>
				</SidebarItem>
			</SidebarSection>
		</SidebarBody>
		<SidebarFooter>
			<SidebarUser />
		</SidebarFooter>
	</Sidebar>
{/snippet}

<div class="relative isolate z-0 flex min-h-svh w-full max-lg:flex-col">
	<!-- Sidebar on desktop -->
	<div class={['fixed inset-y-0 left-0 max-lg:hidden', collapsed ? 'w-[4.5rem]' : 'w-64']}>
		{@render sidebar()}
	</div>

	<!-- Sidebar on mobile -->
	{#if api.open}
		<div
			transition:fade={{ duration: api.open ? 150 : 300 }}
			use:portal
			class="fixed inset-0 bg-black/30"
			{...api.getBackdropProps()}
		></div>
		<div
			transition:fly={{ x: '-100%', duration: 300, opacity: 1 }}
			class="fixed inset-y-0 w-full max-w-80 p-2"
			use:portal
			{...api.getPositionerProps()}
		>
			<div
				id="sidebar"
				class="shadow-xs relative flex h-full flex-col rounded-lg bg-white ring-1 ring-zinc-950/5 dark:bg-zinc-900 dark:ring-white/10"
				{...api.getContentProps()}
			>
				<div class="-mb-3 px-4 pt-3">
					<span class="relative"> </span>
				</div>
				{@render sidebar(true)}
			</div>
		</div>
	{/if}

	<div class="divide-outline-variant bg-surface-low flex flex-col divide-y last:mb-2 lg:hidden">
		<Verify />
		<OrgPlanBanner />
	</div>

	<!-- Navbar on mobile -->
	<header class="flex items-center px-4 lg:hidden">
		<div class="py-2.5">
			<Button plain aria-label="Open navigation" {...api.getTriggerProps() as any}>
				<List data-slot="icon" />
			</Button>
		</div>
		<div class="min-w-0 flex-1 flex-col">
			<!-- Rest of header content	 -->
		</div>
	</header>

	<!-- Content -->
	<main
		class={[
			'z-10 flex flex-1 flex-col pb-2 lg:min-w-0 lg:pr-2 lg:pt-2 lg:transition-[margin]',
			collapsed ? 'lg:ml-[4.5rem]' : 'lg:ml-64'
		]}
	>
		<div class="divide-outline-variant has-last:-mt-2 hidden flex-col divide-y lg:flex">
			<Verify />
			<OrgPlanBanner />
		</div>
		<div
			class="lg:shadow-xs grow px-6 pb-6 lg:rounded-lg lg:bg-white lg:px-10 lg:py-6 lg:pt-6 lg:ring-1 lg:ring-zinc-950/5 dark:lg:bg-zinc-900 dark:lg:ring-white/10"
		>
			{@render children()}
		</div>
	</main>
</div>
