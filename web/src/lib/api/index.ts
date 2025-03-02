import { API_URL } from '$lib/config';
import type { OrganizationAPI } from './organization.ts';
import { getAPI } from 'api-typify';

import type { DefaultError } from '@tanstack/svelte-query';
import type { GitAPI } from './git.ts';
import { ory } from '$lib/queries/auth.js';

export type apiDefs = OrganizationAPI & GitAPI;

export interface CustomProps {
	headers?: Record<string, string>;
	signal?: AbortSignal;
}

export interface APIError extends DefaultError {
	status: number;
	redirectURL?: string;
}

let lastCheck = 0;
let currCheck: Promise<void> | undefined;
function checkSession() {
	// We want to avoid sending off multiple requests to check the session
	// this is quite common when a user focuses a tab and multiple requests are sent
	if (Date.now() - lastCheck < 5_000) return;

	if (currCheck) return currCheck;

	currCheck = ory
		.toSession()
		.catch((error) => {
			if (error?.response?.status === 401) {
				// refresh the page to force a re-login
				globalThis.location.reload();
			}
		})
		.then(() => {
			lastCheck = Date.now();
		});
}

export const createAPI = (extraHeaders: Record<string, string> = {}) =>
	getAPI<apiDefs, CustomProps>(API_URL, (url, options) =>
		fetch(url, {
			...options,
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json',
				...extraHeaders,
				...options?.headers
			},
			credentials: 'include',
			signal: options?.signal
		}).then(async (res) => {
			if (res.ok) {
				return res.json().catch(() => undefined);
			}

			const error = (await res.json().catch(() => res)) as APIError;

			error.status = res.status;

			if (error.redirectURL) {
				globalThis.location.assign(error.redirectURL);
			}

			if (res.status === 401) {
				// If we get a 401, we need to check if the session is still valid
				checkSession();
			}

			return Promise.reject(error);
		})
	);

export const API = createAPI();
export type { Organization } from './organization.ts';
export type { GitInstallation, GitInstallations, GitRepositories, GitRepository } from './git.ts';
