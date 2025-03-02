<script lang="ts">
	import { Skeleton } from '$lib/components';
	import clsx from 'clsx';
</script>

<div class="flex flex-col space-y-2 overflow-hidden">
	{#each Array.from({ length: 10 }) as _, i (i)}
		<Skeleton
			class={clsx('h-4', ['w-1/2', 'w-1/5', 'w-3/4', 'w-2/3'][Math.floor(Math.random() * 4)])}
		/>
	{/each}
</div>
