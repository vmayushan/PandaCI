import { KRATOS_URL } from '$lib/config';
import { Configuration, FrontendApi } from '@ory/client';

export const authAPI = new FrontendApi(
	new Configuration({
		basePath: KRATOS_URL,
		baseOptions: {
			withCredentials: true
		}
	})
);
