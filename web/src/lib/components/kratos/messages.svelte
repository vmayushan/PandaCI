<script lang="ts">
	import type { UiText } from '@ory/client';
	import { Text } from '../text';

	export let messages: UiText[];
</script>

{#if messages && messages.length > 0}
	<ul class="mt-4 flex flex-col space-y-2">
		{#each messages as msg (msg.id)}
			<li>
				{#if msg.type === 'info'}
					<Text>{msg.text}</Text>
				{:else if msg.type === 'success'}
					<Text variant="success">{msg.text}</Text>
				{:else if msg.type === 'error'}
					<Text variant="error">{msg.text}</Text>
				{/if}
			</li>
		{/each}
	</ul>
{/if}
