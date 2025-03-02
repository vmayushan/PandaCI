import { mergeQueryKeys } from '@lukemorales/query-key-factory';
import { organizationQueries } from './organization.js';
import { authQueries } from './auth.js';
import { githubQueries } from './git.js';
import { runQueries } from './runs.js';
import { projectVariableQueries } from './projectVariables.js';
import { projectEnvironmentQueries } from './projectEnvironments.js';

export const queries = mergeQueryKeys(
	organizationQueries,
	authQueries,
	githubQueries,
	runQueries,
	projectVariableQueries,
	projectEnvironmentQueries
);
