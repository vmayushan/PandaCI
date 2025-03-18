import { queries } from '$lib/queries';
import { createQuery } from '@tanstack/svelte-query';
import { goto } from '$app/navigation';

export function getCurrentOrg({
	page
}: {
	page: {
		params: {
			orgName?: string;
		};
		route: {
			id: string | null;
		};
	};
}) {
	// TODO - we should be able to use the orgs data to pre-populate the current org to avoid two requests
	let orgName = $state(page.params.orgName || localStorage.getItem('currentOrg'));

	const orgs = createQuery(() => queries.organization.list());

	const currentOrg = createQuery(() => ({
		...queries.organization.getByName(orgName!),
		enabled: Boolean(orgName)
	}));

	$effect(() => {
		if (page.params.orgName) {
			orgName = page.params.orgName;
		}
	});

	$effect(() => {
		if (orgs.data?.some((o) => o.slug === page.params.orgName)) {
			orgName = page.params.orgName!;
		}
	});

	$effect(() => {
		if (orgName) localStorage.setItem('currentOrg', orgName);
	});

	$effect(() => {
		if (orgs.data && orgs.data.length > 0 && !orgName) {
			orgName = orgs.data[0].slug;
		}
	});

	$effect(() => {
		if (
			orgs.isFetched &&
			// TODO - we should be able to remove this one svelte query properly works with runes
			!orgs.isFetching &&
			!orgs.data?.length &&
			(page.route.id?.startsWith('/[orgName]') || page.route.id === '/')
		) {
			if (localStorage.getItem('welcome-visited') !== 'true')
				goto('/welcome', {
					replaceState: true
				});
			else {
				goto('/account/orgs', {
					replaceState: true
				});
			}
		}
	});

	$effect(() => {
		if (
			currentOrg.isFetched &&
			Boolean(orgs.data?.length) &&
			!orgs.data?.find((org) => org.slug === orgName)
		) {
			orgName = orgs.data![0].slug;
		}
	});

	return currentOrg;
}
