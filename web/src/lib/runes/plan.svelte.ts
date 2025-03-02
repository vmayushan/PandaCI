import { queries } from '$lib/queries';
import { createQuery } from '@tanstack/svelte-query';
import { page } from '$app/state';
import { PUBLIC_STAGE } from '$env/static/public';

const sandboxPrices = {
	PRO: 'pri_01jk6wc23hdjcnzvj68mza3n2k',
	ENTERPRISE: ''
};

const prodPrices = {
	PRO: 'pri_01jkwzdhjk4shw9vkv42gpr5vv',
	ENTERPRISE: ''
};

export const prices = PUBLIC_STAGE === 'prod' ? prodPrices : sandboxPrices;

export function getOrgPlan() {
	const org = createQuery(() => queries.organization.getByName(page.params.orgName));

	const _plan = $derived.by(() => {
		if (!org.data) return undefined;
		return org.data.license?.plan || 'paused';
	});

	return {
		get plan() {
			return _plan;
		},
		get features() {
			return (
				org.data?.license?.features || {
					maxProjects: 10,
					maxCommitters: 5,
					committers: 5,
					maxCloudRunnersScale: 4,
					maxBuidMinutes: 6000,
					buildMinutes: 0,
					buildMinutesPriceID: ''
				}
			);
		},
		get status() {
			return org.data?.license?.paddleData?.subscriptionStatus;
		}
	};
}
