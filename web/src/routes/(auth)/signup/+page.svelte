<script lang="ts">
	import { page } from '$app/state';
	import { Card } from '$lib/components';
	import Divider from '$lib/components/divider.svelte';
	import Heading from '$lib/components/heading.svelte';
	import Messages from '$lib/components/kratos/messages.svelte';
	import Render from '$lib/components/kratos/render.svelte';
	import { Text, TextLink } from '$lib/components/text';
	import { authAPI } from '$lib/kratos';
	import type { SvelteHTMLElements } from 'svelte/elements';
	import { handleError } from '../handleError.svelte';
	import { replaceState } from '$app/navigation';
	import { ArrowLeft } from 'phosphor-svelte';
	import Button from '$lib/components/button.svelte';
	import Skeleton from '$lib/components/skeleton.svelte';

	const flowId = page.url.searchParams.get('flow');
	const returnTo = decodeURIComponent(page.url.searchParams.get('return_to') || '') || undefined;

	const session = flowId
		? (authAPI
				.getRegistrationFlow({
					id: flowId
				})
				.then((res) => {
					page.url.searchParams.set('flow', res.data.id);
					replaceState(page.url, page.state);
					return res;
				})
				.catch(handleError) as ReturnType<typeof authAPI.getRegistrationFlow>)
		: (authAPI
				.createBrowserRegistrationFlow({ returnTo })
				.then((res) => {
					page.url.searchParams.set('flow', res.data.id);
					replaceState(page.url, page.state);
					return res;
				})
				.catch(handleError) as ReturnType<typeof authAPI.createBrowserRegistrationFlow>);
</script>

<div class="flex flex-1 grow flex-col justify-center px-2">
	{#if session}
		{#await session}
			<Card class="mx-auto my-14 flex w-full max-w-md flex-col">
				<div>
					<Heading size="sm" level={2}>Sign up</Heading>
					<Skeleton class="mt-2 h-6 w-3/5" />
				</div>
				<Skeleton class="mt-12 h-11 w-full sm:h-9" />

				<Text class="mt-8">
					By signing up, you agree to our <TextLink
						href="https://pandaci.com/legal"
						target="_blank"
					>
						Terms and Conditions
					</TextLink>
				</Text>
			</Card>
		{:then { data }}
			<Button
				href="https://pandaci.com"
				class="!absolute left-0 top-2 w-min sm:left-8 sm:top-8"
				plain
			>
				<ArrowLeft data-slot="icon" />
				Home
			</Button>
			<Card class="mx-auto my-14 flex w-full max-w-md flex-col space-y-12 ">
				<div>
					<Heading size="sm" level={2}>Sign up</Heading>
					<Text class="mt-2">
						Already have an account? <TextLink
							href={`/login${data.return_to ? `?return_to=${encodeURIComponent(data.return_to)}` : ''}`}
						>
							Login
						</TextLink>
					</Text>
				</div>

				{#if data.ui.messages}
					<Messages messages={data.ui.messages} />
					<Divider class="my-2" />
				{/if}
				<form
					class="flex w-full flex-col"
					method={data.ui.method as SvelteHTMLElements['form']['method']}
					action={data.ui.action}
				>
					<div class="flex flex-col space-y-8">
						<Render nodes={data.ui.nodes} filters={{ groups: ['oidc'] }} />
					</div>
					<Text class="mt-8">
						By signing up, you agree to our <TextLink
							href="https://pandaci.com/legal"
							target="_blank"
						>
							Terms and Conditions
						</TextLink>
					</Text>
				</form>
			</Card>
		{/await}
	{/if}
</div>
