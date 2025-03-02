<script lang="ts">
	import { page } from '$app/state';
	import { authAPI } from '$lib/kratos';
	import { Card } from '$lib/components';
	import Messages from '$lib/components/kratos/messages.svelte';
	import Render from '$lib/components/kratos/render.svelte';
	import Heading from '$lib/components/heading.svelte';
	import { handleError } from '../handleError.svelte';
	import { replaceState } from '$app/navigation';
	import Fieldset from '$lib/components/fieldset/fieldset.svelte';
	import FieldGroup from '$lib/components/fieldset/fieldGroup.svelte';

	const flowId = page.url.searchParams.get('flow');

	const email = page.url.searchParams.get('email');

	const session = flowId
		? authAPI
				.getVerificationFlow({
					id: flowId
				})
				.catch(handleError)
		: authAPI
				.createBrowserVerificationFlow({
					returnTo: page.url.searchParams.get('return_to') || undefined
				})
				.then(async (res) => {
					page.url.searchParams.set('flow', res.data.id);
					replaceState(page.url, page.state);
					return res;
				})
				.catch(handleError);
</script>

{#if session}
	<Card class="mx-auto my-auto flex w-full max-w-md flex-col space-y-12 ">
		<div>
			<Heading>Verification</Heading>

			{#await session then data}
				{#if data}
					{#if data.data.ui.messages}
						<Messages messages={data.data.ui.messages} />
					{/if}
				{/if}
			{/await}
		</div>

		{#await session then data}
			{#if data}
				<form class="w-full" method={data.data.ui.method as any} action={data.data.ui.action}>
					<Fieldset>
						<FieldGroup>
							<Render
								defaultValues={email ? { email } : {}}
								filters={{
									excludeAttributes: [{ name: 'email', type: 'submit' }, { type: 'hidden' }]
								}}
								nodes={data.data.ui.nodes}
							/>
						</FieldGroup>
					</Fieldset>
					<Render
						defaultValues={email ? { email } : {}}
						filters={{
							attributes: [{ type: 'hidden' }]
						}}
						nodes={data.data.ui.nodes}
					/>
				</form>
			{/if}
		{/await}
	</Card>
{/if}
