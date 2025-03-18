<script lang="ts">
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { ArrowRight } from 'phosphor-svelte';
	import { page } from '$app/state';

	interface VerifyProps {
		class?: string;
	}

	const org = createQuery(() => ({
		...queries.organization.getByName(page.params.orgName),
		enabled: Boolean(page.params.orgName)
	}));

	const usage = createQuery(() => ({
		...queries.organization.getByName(page.params.orgName)._ctx.usage(),
		enabled: Boolean(page.params.orgName)
	}));

	const { class: className }: VerifyProps = $props();
</script>

{#if page.params.orgName}
	{#if org.data?.license?.plan === 'paused'}
		<div class={['flex w-full items-center justify-center gap-x-6 px-6 py-2.5', className]}>
			<a
				class="text-on-surface-variant hover:text-on-surface flex items-center justify-center space-x-1 text-sm"
				href={`/${page.params.orgName}/~/billing`}
			>
				<span>A paid plan is required for this org</span>
				<ArrowRight data-slot="icon" />
			</a>
		</div>
	{:else if org.data?.license?.plan === 'free' && (usage.data?.usedBuildMinutes || 0) >= org.data?.license.features.maxBuidMinutes}
		<div class={['flex w-full items-center justify-center gap-x-6 px-6 py-2.5', className]}>
			<a
				class="text-on-surface-variant hover:text-on-surface flex items-center justify-center space-x-1 text-sm"
				href={`/${page.params.orgName}/~/billing`}
			>
				<span>Monthly free build minutes exceeded, please upgrade to continue building</span>
				<ArrowRight class="shrink-0" data-slot="icon" />
			</a>
		</div>
	{:else if org.data?.license?.plan === 'free' && (usage.data?.usedCommitters ?? 0) >= org.data?.license.features.maxCommitters}
		<div class={['flex w-full items-center justify-center gap-x-6 px-6 py-2.5', className]}>
			<a
				class="text-on-surface-variant hover:text-on-surface flex items-center justify-center space-x-1 text-sm"
				href={`/${page.params.orgName}/~/billing`}
			>
				<span>
					You've reached the maximum number of committers allowed on the free plan, please upgrade
					if you need more</span
				>
				<ArrowRight class="shrink-0" data-slot="icon" />
			</a>
		</div>
	{/if}
{/if}
