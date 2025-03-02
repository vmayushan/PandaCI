import { queries } from '$lib/queries';
import { createQuery } from '@tanstack/svelte-query';
import { goto } from '$app/navigation';

export interface User {
	id: string;
	name: string;
	email: string;
	avatar?: string;
}

export function getUser({
	enabled
}: {
	enabled?: boolean;
} = {}) {
	return createQuery(() => ({
		...queries.auth.session(),
		select: ({ data }) => ({ id: data.identity?.id, ...data.identity?.traits }) as User,
		retry(failureCount, error: any) {
			if (error?.response?.status === 401) {
				if (globalThis.location.pathname === '/') goto('/login');
				else goto(`/login?return_to=${encodeURIComponent(globalThis.location.toString())}`);
				return false;
			}
			return failureCount < 2;
		},
		enabled: enabled ?? true
	}));
}
