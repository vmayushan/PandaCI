<script lang="ts">
	import type { Snippet } from 'svelte';

	interface ContainerLayoutProps {
		children: Snippet;
	}

	const { children }: ContainerLayoutProps = $props();
</script>

<div class="mx-auto max-w-7xl">
	{@render children()}
</div>
