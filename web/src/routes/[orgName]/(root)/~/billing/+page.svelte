<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { API } from '$lib/api';
	import { Subheading } from '$lib/components';
	import Button from '$lib/components/button.svelte';
	import Card from '$lib/components/card.svelte';
	import DescriptionDetails from '$lib/components/descriptionList/DescriptionDetails.svelte';
	import DescriptionList from '$lib/components/descriptionList/DescriptionList.svelte';
	import DescriptionTerm from '$lib/components/descriptionList/DescriptionTerm.svelte';
	import { type PaddleContext } from '$lib/components/paddle/paddleContext.svelte';
	import SubHeading from '$lib/components/subHeading.svelte';
	import Text from '$lib/components/text/text.svelte';
	import TextLink from '$lib/components/text/textLink.svelte';
	import Title from '$lib/components/title.svelte';
	import { queries } from '$lib/queries';
	import { getOrgPlan } from '$lib/runes/plan.svelte';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';

	const { openCheckout } = getContext<PaddleContext>('paddle');

	const getPortalURL = createMutation(() => ({
		mutationFn: () =>
			API.post('/v1/paddle/orgs/{orgSlug}/portal', { params: { orgSlug: page.params.orgName } })
	}));

	const usage = createQuery(() => queries.organization.getByName(page.params.orgName)._ctx.usage());

	const queryClient = useQueryClient();

	$effect(() => {
		if (page.url.searchParams.get('s') !== null) {
			queryClient
				.invalidateQueries(queries.organization.getByName(page.params.orgName))
				.then(() => {
					goto(`/${page.params.orgName}/~/billing`);
				});
		}
	});

	const plan = getOrgPlan();

	const isLoading = $derived(Boolean(page.url.searchParams.get('s')));
</script>

<Title title="Billing">
	{#snippet description()}
		<TextLink href="https://pandaci.com/pricing">Learn more</TextLink> about our pricing plans and usage.
		If you have any questions, please contact us at <TextLink href="mailto:support@pandaci.com">
			support@pandaci.com
		</TextLink>
	{/snippet}
</Title>

{#if plan.plan === 'pro' && plan.status !== 'canceled' && !isLoading}
	<Card class="mt-10">
		<Subheading>Manage your subscription</Subheading>
		<Text>You are currently subscribed to the Pro plan.</Text>
		<Button
			class="mt-4"
			loading={getPortalURL.isPending}
			onclick={() =>
				getPortalURL.mutate(undefined, {
					onSuccess(data) {
						window.open(data.generalURL, '_blank');
					}
				})}
		>
			Manage
		</Button>
	</Card>
{/if}

{#if (['free', 'paused'].includes(plan.plan || '') || plan.status === 'canceled') && !isLoading}
	<section class="mx-auto mt-12">
		<div
			class="mx-auto grid max-w-lg grid-cols-1 items-center gap-y-6 sm:gap-y-0 lg:max-w-4xl lg:grid-cols-2"
		>
			<div
				class="bg-surface ring-outline-variant rounded-3xl rounded-t-3xl p-8 ring-1 sm:mx-8 sm:rounded-b-none sm:p-10 lg:mx-0 lg:rounded-bl-3xl lg:rounded-tr-none"
			>
				<h3 id="tier-hobby" class="text-base/7 font-semibold text-green-500">Pro</h3>
				<p class="mt-4 flex items-baseline gap-x-2">
					<span class="text-on-surface text-3xl font-semibold tracking-tight xl:text-5xl">$10</span>
					<span class="text-on-surface-variant text-base">/month</span>
				</p>
				<p class="text-on-surface-variant mt-6 text-base/7">
					The perfect plan if you&#039;re just getting started with our product.
				</p>
				<ul role="list" class="text-on-surface-variant mt-8 space-y-3 text-sm/6 sm:mt-10">
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Up to 100 projects
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Up to 8 core runners
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Includes 6000 build minutes
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						$1 per 500 additional build minutes
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Email support
					</li>
				</ul>
				<Button onclick={() => openCheckout()} class="mt-8" full color="green">Upgrade</Button>
			</div>
			<div class="bg-surface-high ring-outline relative rounded-3xl p-8 shadow-2xl ring-1 sm:p-10">
				<h3 id="tier-enterprise" class="text-base/7 font-semibold text-green-500">Enterprise</h3>
				<p class="mt-4 flex items-baseline gap-x-2">
					<span class="text-on-surface text-3xl font-semibold tracking-tight xl:text-5xl">
						From $600
					</span>
					<span class="text-on-surface-variant text-base">/month</span>
				</p>
				<p class="text-on-surface-variant mt-6 text-base/7">
					Dedicated support and infrastructure for your company.
				</p>
				<ul role="list" class="text-on-surface-variant mt-8 space-y-3 text-sm/6 sm:mt-10">
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Unlimited projects
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Up to 16 core runners
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Volume discounts
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Dedicated support representative
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						SAML (coming soon)
					</li>
					<li class="flex gap-x-3">
						<svg
							class="h-6 w-5 flex-none text-green-600"
							viewBox="0 0 20 20"
							fill="currentColor"
							aria-hidden="true"
							data-slot="icon"
						>
							<path
								fill-rule="evenodd"
								d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
								clip-rule="evenodd"
							/>
						</svg>
						Custom domain
					</li>
				</ul>
				<Button class="mt-8" color="green" full href="mailto:support@pandaci.com">
					Chat with us
				</Button>
			</div>
		</div>
	</section>
{/if}

{#if !isLoading}
	<section class="my-12">
		<Card>
			<SubHeading>Usage</SubHeading>
			<DescriptionList>
				<DescriptionTerm>Build minutes</DescriptionTerm>
				<DescriptionDetails>
					{usage.data?.usedBuildMinutes} / {plan.features.buildMinutes}
				</DescriptionDetails>
				<DescriptionTerm>Commiters</DescriptionTerm>
				<DescriptionDetails>
					{usage.data?.usedCommitters} / {plan.features.committers}
				</DescriptionDetails>
				<DescriptionTerm>Projects</DescriptionTerm>
				<DescriptionDetails>
					{usage.data?.projectCount} / {plan.features.maxProjects}
				</DescriptionDetails>
			</DescriptionList>
			<Text class="mt-4">
				If you're on a Pro or Enterprise plan, we'll automatically increase your build minutes and
				committers once you reach the limit. <br /> Build minutes are charged in increments of 500 minutes.
			</Text>
		</Card>
	</section>
{/if}
