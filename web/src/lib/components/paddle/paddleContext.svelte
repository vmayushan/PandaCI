<script lang="ts" module>
	export type Plan = 'pro' | 'enterprise' | 'free' | 'paused';

	export interface PaddleContext {
		openCheckout: () => void;
	}

	export function useCheckout() {
		return getContext<PaddleContext>('paddle');
	}
</script>

<script lang="ts">
	import { initializePaddle } from '@paddle/paddle-js';
	import { PUBLIC_STAGE } from '$env/static/public';
	import { getTheme } from '../themeProvider.svelte';
	import { getContext, setContext, type Snippet } from 'svelte';
	import { queries } from '$lib/queries';
	import { createQuery } from '@tanstack/svelte-query';
	import { page } from '$app/state';
	import { prices } from '$lib/runes/plan.svelte';

	interface PaddleProps {
		children: Snippet;
	}

	const { children }: PaddleProps = $props();

	const paddlePromise = initializePaddle({
		environment: PUBLIC_STAGE === 'prod' ? 'production' : 'sandbox',
		token:
			PUBLIC_STAGE === 'prod'
				? 'live_beb62518ab63d6938d987a717c9'
				: 'test_b38e5ee92f56be3fe2874e2ddbb'
	});

	const theme = getTheme();

	function getPaddleCtx() {
		const org = createQuery(() => queries.organization.getByName(page.params.orgName));

		return {
			async openCheckout() {
				const paddle = await paddlePromise;
				paddle?.Checkout.open({
					settings: {
						theme: theme.resolvedTheme(),
						displayMode: 'overlay',
						variant: 'one-page',
						allowLogout: true,
						successUrl: `https://app.pandaci.com/${org.data?.slug}/~/billing?s=success`
					},
					customData: {
						orgID: org.data!.id
					},
					items: [
						{
							priceId: prices.PRO,
							quantity: 1
						}
					]
				});
			}
		};
	}

	const ctx = getPaddleCtx();

	setContext<PaddleContext>('paddle', ctx);
</script>

{@render children()}
