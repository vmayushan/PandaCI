import { KRATOS_URL } from '$lib/config';
import { createQueryKeys } from '@lukemorales/query-key-factory';
import { Configuration, FrontendApi } from '@ory/client';

export const ory = new FrontendApi(
	new Configuration({
		basePath: KRATOS_URL,
		baseOptions: {
			withCredentials: true
		}
	})
);

export const authQueries = createQueryKeys('auth', {
	session: () => ({
		queryFn: () => ory.toSession(),
		queryKey: ['session']
	})
});
