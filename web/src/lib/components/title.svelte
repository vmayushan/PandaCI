<script lang="ts">
	import type { Snippet } from 'svelte';
	import Heading from './heading.svelte';
	import type {
		RunStatus as RunStatusType,
		RunConclusion as RunConclusionType
	} from '$lib/api/organization';
	import RunStatus from './runStatus.svelte';
	import Text from './text/text.svelte';
	import Skeleton from './skeleton.svelte';

	interface TitleProps {
		title: string;
		action?: Snippet;
		description?: Snippet;
		children?: Snippet;
		state?: {
			status?: RunStatusType;
			conclusion?: RunConclusionType;
		};
		titleLoading?: boolean;
		class?: string;
	}

	const {
		title,
		action,
		class: className,
		titleLoading,
		children,
		state,
		description
	}: TitleProps = $props();
</script>

<header
	class={[
		'border-outline-variant flex min-h-16 w-full flex-wrap items-center justify-between gap-4 border-b pb-5',
		!action && 'mt-4',
		className
	]}
>
	<div class="flex min-w-0 grow flex-col">
		{@render action?.()}
		<Heading class={['flex max-w-full items-center space-x-4', action && 'mt-4']}>
			{#if state}
				<RunStatus
					class="mb-auto mt-1 size-6"
					status={state.status}
					conclusion={state.conclusion}
				/>
			{/if}
			{#if titleLoading}
				<Skeleton class="h-6 w-32" />
			{:else}
				<span class="min-w-0 grow hyphens-auto break-words">
					{title}
				</span>
			{/if}
		</Heading>
		{#if description}
			<Text>
				{@render description()}
			</Text>
		{/if}
	</div>
	<div class="flex gap-4">
		{#if children}
			{@render children()}
		{/if}
	</div>
</header>
