<script lang="ts">
	import { page } from '$app/state';
	import { authAPI } from '$lib/kratos';
	import Messages from '$lib/components/kratos/messages.svelte';
	import Render from '$lib/components/kratos/render.svelte';
	import Heading from '$lib/components/heading.svelte';
	import { Text, TextLink } from '$lib/components/text';
	import { Card } from '$lib/components';
	import { goto, replaceState } from '$app/navigation';
	import { handleForm } from '$lib/utils';
	import type { UpdateLoginFlowBody } from '@ory/client';
	import { handleError } from '../handleError.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { queries } from '$lib/queries';
	import Button from '$lib/components/button.svelte';
	import { ArrowLeft } from 'phosphor-svelte';

	const flowId = page.url.searchParams.get('flow');
	const returnTo = decodeURIComponent(page.url.searchParams.get('return_to') || '') || undefined;
	const refresh = page.url.searchParams.get('refresh') === 'true';

	const session = flowId
		? authAPI
				.getLoginFlow({
					id: flowId
				})
				.catch(handleError)
		: authAPI
				.createBrowserLoginFlow({ returnTo, refresh })
				.then(async (res) => {
					page.url.searchParams.set('flow', res.data.id);
					replaceState(page.url, page.state);
					return res;
				})
				.catch(handleError);

	const queryClient = useQueryClient();
</script>

{#await session then sessionData}
	{#if sessionData}
		{@const data = sessionData.data}
		<div class="flex flex-1 grow flex-col justify-center px-2">
			{#if !sessionData.data.refresh}
				<Button
					href="https://pandaci.com"
					class="!absolute left-0 top-2 w-min sm:left-8 sm:top-8"
					plain
				>
					<ArrowLeft data-slot="icon" />
					Home
				</Button>
			{/if}
			<Card class="mx-auto my-14 flex w-full max-w-md flex-col space-y-12">
				<div>
					<Heading size="sm" level={2}>Login</Heading>
					{#if !data.ui.messages?.some((m) => m.id === 1010003)}
						<Text class="mt-2">
							Dont have an account? <TextLink
								href={`/signup${data.return_to ? `?return_to=${encodeURIComponent(data.return_to)}` : ''}`}
							>
								Sign up for free
							</TextLink>
						</Text>
					{/if}

					{#if data.ui.messages}
						<Messages messages={data.ui.messages} />
					{/if}
				</div>

				<form
					class="flex w-full flex-col"
					onsubmit={(e) => {
						const { data: body, value: method } = handleForm<UpdateLoginFlowBody>(e);

						authAPI
							.updateLoginFlow({
								flow: data.id,
								updateLoginFlowBody: { ...body, method: 'oidc', provider: method as any }
							})
							.then((res) => {
								queryClient.setQueryData(queries.auth.session().queryKey, res.data.session);
								queryClient.prefetchQuery(queries.organization.list());
								goto(data.return_to || returnTo || '/');
							})
							.catch(handleError);
					}}
				>
					<Render nodes={data.ui.nodes} filters={{ groups: ['default'] }} />
					<div class="flex flex-col space-y-8">
						<Render nodes={data.ui.nodes} filters={{ groups: ['oidc'] }} />
					</div>

					<Text class="mt-8">
						By signing in, you agree to our <TextLink
							href="https://pandaci.com/legal"
							target="_blank"
						>
							Terms and Conditions
						</TextLink>
					</Text>
				</form>
			</Card>
		</div>
	{/if}
{/await}
