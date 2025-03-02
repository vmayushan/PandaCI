import { $, docker, env, job } from "@pandaci/workflow";

async function initGo() {
  await $`go mod download`;
  await $`go mod verify`;
}

async function initFly() {
  await $`apt-get update && apt-get install -y curl`;
  await $`curl -L https://fly.io/install.sh | sh`;
  await $`ln -s /root/.fly/bin/flyctl /usr/local/bin/flyctl`;
}

export const runBackendStage = (stage: "dev" | "prod") =>
  job.nothrow(`Backend - ${stage}`, async () => {
    await docker.if(stage === "dev")(
      "golang:1.24.0",
      { name: "Test" },
      async () => {
        await initGo();
        await $`go test ./...`;
      },
    );

    await docker.if(stage === "dev")(
      "node:22",
      { name: "Deploy - preview" },
      async () => {
        await initFly();

        const flyEnv = {
          FLY_ACCESS_TOKEN: env.FLY_ACCESS_TOKEN_DEV,
          FLY_APP:
            `pandaci-core-prev-${env.PANDACI_BRANCH}-${env.PANDACI_PR_NUMBER}`,
        };

        // if app doesn't exist, create it
        if (
          (await $`flyctl status`.env(flyEnv)
            .nothrow()).exitCode !== 0
        ) {
          await $`flyctl apps create ${flyEnv.FLY_APP} --org pandaci-dev`.env(
            flyEnv,
          );

          await $`flyctl secrets set DEV_BRANCH=${env.PANDACI_BRANCH}`.env(
            flyEnv,
          );

          await $`npx wrangler kv:key put ${env.PANDACI_BRANCH} ${flyEnv.FLY_APP} --namespace-id=${env.CLOUDFLARE_API_KV_ID}`
            .env({
              CLOUDFLARE_API_TOKEN: env.CLOUDFLARE_API_TOKEN,
              CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
            });
        }

        await $`echo $INPUT_SECRETS | tr " " "\n" | flyctl secrets import -a ${flyEnv.FLY_APP}`
          .env({
            ...flyEnv,
            INPUT_SECRETS: env.FLY_APP_SECRETS,
          });

        await $`flyctl -c ./ee/fly/core/fly.toml deploy`.env(flyEnv);
      },
    );

    await docker.if(stage === "prod")(
      "ubuntu:24.04",
      { name: "Deploy - main" },
      async () => {
        await initFly();

        await $`flyctl -c ./ee/fly/core/fly.toml deploy`.env({
          FLY_ACCESS_TOKEN: env.FLY_ACCESS_TOKEN_PROD,
          FLY_APP: "pandaci-core-prod",
        });
      },
    );
  });
