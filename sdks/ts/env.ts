import { DEFAULTS } from "./defaults.ts";

type EnvVariables = {
  PANDACI_WORKFLOW_ID: string;
  PANDACI_REPO_URL: string;
  PANDACI_BRANCH: string;
  PANDACI_PR_NUMBER?: string;
  PANDACI_COMMIT_SHA: string;
  PANDACI_DENO_VERSION: string;
};

type InternalEnvVariables = {
  PANDACI_WORKFLOW_GRPC_ADDRESS: string;
  PANDACI_WORKFLOW_JWT: string;
  PANDACI_LOG_LEVEL: number;
};

/**
 * PandaCI environment variables type. Add a generic to have typesafe envs
 * @example const safeEnv = env as Env<{test: string}>
 */
export type Env<
  T extends Record<string | number | symbol, string> = Record<
    string | number | symbol,
    string
  >,
> =
  & {
    [key: `PANDACI_${string}`]: never;
  }
  & T
  & EnvVariables;

export type ExtendedEnv = Env & InternalEnvVariables;

const combinedEnvVariables: ExtendedEnv = {
  PANDACI_LOG_LEVEL: DEFAULTS.logLevel,
  ...Deno.env.toObject(),
} as ExtendedEnv;

/**
 * PandaCI env variables. These are a combination of predefined variables and
 * ones you have declared in the secrets tab in the dashboard
 */
export const env: Env = combinedEnvVariables;
export const extendedEnv = combinedEnvVariables;
