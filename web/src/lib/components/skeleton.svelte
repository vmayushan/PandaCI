<script lang="ts">
	import clsx from 'clsx';

	interface SkeletonProps {
		class?: string;
		loading?: boolean;
	}

	const { class: className, loading }: SkeletonProps = $props();
</script>

{#if loading ?? true}
	<div class={clsx(className, 'animate-pulse rounded-md bg-black/10 dark:bg-white/5')}></div>
{/if}
