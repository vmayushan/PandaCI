<script lang="ts">
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { clsx } from 'clsx';
	import { Info, Warning, XCircle } from 'phosphor-svelte';

	type AlertProps = {
		type: 'info' | 'warning' | 'error';
		title: string;
		message: string;
	} & Omit<SvelteHTMLElements['div'], 'children'>;

	const { class: className, type, title, message, ...props }: AlertProps = $props();
</script>

<div
	class={clsx(
		'rounded-md border p-4',
		type === 'info' && 'border-blue-200 bg-blue-50 dark:border-blue-300/25 dark:bg-blue-900/25',
		type === 'warning' &&
			'border-amber-200 bg-amber-50 dark:border-amber-300/25 dark:bg-amber-900/25',
		type === 'error' && 'border-red-200 bg-red-50 dark:border-red-300/25 dark:bg-red-900/25',
		className
	)}
>
	<div class="flex">
		<div class="shrink-0">
			{#if type === 'error'}
				<XCircle data-slot="icon" class="size-5 text-red-400 dark:text-red-100" />
			{:else if type === 'info'}
				<Info data-slot="icon" class="size-5 text-blue-400 dark:text-blue-100" />
			{:else if type === 'warning'}
				<Warning data-slot="icon" class="size-5 text-amber-400 dark:text-amber-100" />
			{/if}
		</div>
		<div class="ml-3">
			<h3
				class={clsx(
					'text-sm font-medium dark:text-white',

					type === 'info' && 'text-blue-800',
					type === 'warning' && 'text-amber-800',
					type === 'error' && 'text-red-800'
				)}
			>
				{title}
			</h3>
			<div
				class={clsx(
					'mt-2 text-sm',

					type === 'info' && 'text-blue-700 dark:text-blue-100',
					type === 'warning' && 'text-amber-700 dark:text-amber-100',
					type === 'error' && 'text-red-700 dark:text-red-100'
				)}
			>
				<p>
					{message}
				</p>
			</div>
		</div>
	</div>
</div>
