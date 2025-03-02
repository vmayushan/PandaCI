import { goto } from '$app/navigation';
import type { AxiosError } from 'axios';

export const handleError = async (error: AxiosError<any, unknown>): Promise<void> => {
	if (!error.response || error.response?.status === 0) {
		window.location.href = `/auth-error?error=${encodeURIComponent(JSON.stringify(error.response))}`;
		return;
	}

	const responseData = error.response?.data || {};

	switch (error.response.status) {
		case 400: {
			if (responseData.error?.id === 'session_already_available') {
				return goto('/');
			}

			console.error('400 error', responseData);
			throw error;
		}
		case 401: {
			// We have no session or the session is invalid
			return goto('/login');
		}
		case 403: {
			if (
				responseData.error?.id === 'session_refresh_required' &&
				responseData.redirect_browser_to
			) {
				return (window.location = responseData.redirect_browser_to);
			}
			break;
		}
		case 404: {
			const errorMsg = {
				data: error.response?.data || error,
				status: error.response?.status,
				statusText: error.response?.statusText,
				url: window.location.href
			};

			return goto(`/auth-error?error=${encodeURIComponent(JSON.stringify(errorMsg))}`);
		}
		case 410: {
			if (responseData.use_flow_id !== undefined) {
				console.warn('sdkError 410: Update flow');
				return goto(`?flow=${responseData.use_flow_id}`);
			}

			return goto(globalThis.location.pathname);
		}
		// we need to parse the response and follow the `redirect_browser_to` URL
		// this could be when the user needs to perform a 2FA challenge
		// or passwordless login
		case 422: {
			if (responseData.redirect_browser_to !== undefined) {
				const currentUrl = new URL(window.location.href);
				const redirect = new URL(responseData.redirect_browser_to);

				// host name has changed, then change location
				if (currentUrl.host !== redirect.host) {
					console.warn('sdkError 422: Host changed redirect');
					window.location = responseData.redirect_browser_to;
					return;
				}

				// Path has changed
				if (currentUrl.pathname !== redirect.pathname) {
					console.warn('sdkError 422: Path changed redirect');
					return goto(redirect.pathname + redirect.search);
				}

				// for webauthn we need to reload the flow
				const flowId = redirect.searchParams.get('flow');

				if (flowId != null) {
					// get new flow data based on the flow id in the redirect url
					console.warn('sdkError 422: Update flow');
					return goto(globalThis.location.pathname + `?flow=${flowId}`);
				}
				console.warn('sdkError 422: Redirect browser to');
				window.location = responseData.redirect_browser_to;
				return;
			}
		}
	}

	const params = new URLSearchParams();
	params.set('error', JSON.stringify(error.response, null, 2));
	params.set('id', responseData.error?.id || '');

	return goto(`/auth-error?${params.toString()}`);
};
