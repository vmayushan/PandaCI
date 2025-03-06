// Note: This is normally "jsr:@pandaci/workflow", but we're just using the local source here.
import { $, docker } from "@pandaci/workflow";
import { env, initFly } from "../utils.ts";

export function testGoTask() {
  return docker("golang:1.24.1", { name: "Testing Go code" }, async () => {
    await $`go mod download`;
    await $`go mod verify`;
    await $`go test ./...`;
  });
}

export function deployDevCoreBackendTask() {
  return docker("node:22", { name: "Deploying core backend" }, async () => {
    await initFly();

    // Setup environment variables
    const flyEnv = {
      FLY_ACCESS_TOKEN: env.FLY_ACCESS_TOKEN_DEV,
      FLY_APP: env.PANDACI_PR_NUMBER
        ? `pandaci-core-prev-${env.PANDACI_BRANCH}-${env.PANDACI_PR_NUMBER}`
        : `pandaci-core-prev-${env.PANDACI_BRANCH}`,
    };

    // Create new app if it doesn't exist
    const appExists =
      (await $`flyctl status`.env(flyEnv).nothrow()).exitCode === 0;
    if (!appExists) {
      await $`flyctl apps create ${flyEnv.FLY_APP} --org pandaci-dev`.env(
        flyEnv,
      );

      await $`flyctl secrets set DEV_BRANCH=${env.PANDACI_BRANCH}`.env(flyEnv);

      // We store the fly app name in cloudflare KV store for our dev proxy worker
      // which allows us to route traffic to the correct deployment
      const cloudflareEnv = {
        CLOUDFLARE_API_TOKEN: env.CLOUDFLARE_API_TOKEN,
        CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
      };

      await $`npx wrangler kv:key put ${env.PANDACI_BRANCH} ${flyEnv.FLY_APP} --namespace-id=${env.CLOUDFLARE_API_KV_ID}`
        .env(cloudflareEnv);
    }

    await $`echo $INPUT_SECRETS | tr " " "\n" | flyctl secrets import -a ${flyEnv.FLY_APP}`
      .env({
        ...flyEnv,
        INPUT_SECRETS: env.FLY_APP_SECRETS,
      });

    await $`flyctl -c ./ee/fly/core/fly.toml deploy`.env(flyEnv);
  });
}

export function deployProdCoreBackendTask() {
  return docker("node:22", { name: "Deploying core backend" }, async () => {
    await initFly();

    await $`flyctl -c ./ee/fly/core/fly.toml deploy`.env({
      FLY_ACCESS_TOKEN: env.FLY_ACCESS_TOKEN_PROD,
      FLY_APP: "pandaci-core-prod",
    });
  });
}
