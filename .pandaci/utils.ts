// Note: This is normally "jsr:@pandaci/workflow", but we're just using the local source here.
import { $, type Env, env as rawEnv } from "@pandaci/workflow";

export const env = rawEnv as Env<{
  PANDACI_BRANCH: string;
  FLY_ACCESS_TOKEN_DEV: string;
  CLOUDFLARE_API_TOKEN: string;
  CLOUDFLARE_ACCOUNT_ID: string;
  CLOUDFLARE_KV_NAMESPACE_ID: string;
  CLOUDFLARE_API_KV_ID: string;
  FLY_APP_SECRETS: string;
  FLY_ACCESS_TOKEN_PROD: string;
}>;

export async function initFly() {
  await $`apt-get update && apt-get install -y curl`;
  await $`curl -L https://fly.io/install.sh | sh`;
  await $`ln -s /root/.fly/bin/flyctl /usr/local/bin/flyctl`;
}

export async function initPnpm() {
  await $`PNPM_HOME="/pnpm"`;
  await $`PATH="$PNPM_HOME:$PATH"`;
  await $`corepack enable`;
  await $`pnpm i`;
}
