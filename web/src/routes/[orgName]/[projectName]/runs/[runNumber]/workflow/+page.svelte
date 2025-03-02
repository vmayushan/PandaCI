<script lang="ts">
	import { page } from '$app/state';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import LogSkeleton from './logSkeleton.svelte';
	import { Logs } from './logs.svelte';
	import RunLogs from '../runLogs.svelte';
	import { Heading } from '$lib/components';
	import Text from '$lib/components/text/text.svelte';

	const run = createQuery(() =>
		queries.runs
			.projectRuns(page.params.orgName, page.params.projectName)
			._ctx.get(page.params.runNumber)
	);

	const logMeta = createQuery(() => ({
		...queries.runs
			.projectRuns(page.params.orgName, page.params.projectName)
			._ctx.logs(run.data?.id ?? '', 'workflow'),
		enabled: !!run.data,
		staleTime: Infinity
	}));

	const rawLogs = new Logs(page.params.orgName, page.params.projectName, 'jsonl');

	$effect(() => {
		rawLogs.update({
			logURL: logMeta.data?.url,
			runID: run.data?.id,
			status: run.data?.status
		});
	});
</script>

<div
	class="border-outline-variant bg-surface-low w-full grow overflow-auto rounded-lg border shadow"
>
	<header class="border-outline-variant bg-surface flex flex-col justify-center border-b px-4 py-2">
		<Heading level={3} size="xs" class="flex items-center space-x-2">
			<span>Workflow logs</span>
		</Heading>
		<Text>
			View detailed workflow execution logs here to identify and troubleshoot any issues during the
			run.
		</Text>
	</header>
	<div class="p-4">
		{#if rawLogs.isLoading || logMeta.isLoading || run.isLoading}
			<LogSkeleton />
		{:else}
			<RunLogs>
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html rawLogs.htmlLogs}
			</RunLogs>
		{/if}
	</div>
</div>
