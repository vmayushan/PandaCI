import { $, docker, env, job } from "@pandaci/workflow";

export async function initPnpm() {
  await $`PNPM_HOME="/pnpm"`;
  await $`PATH="$PNPM_HOME:$PATH"`;
  await $`corepack enable`;
  await $`SHA_SUM=$(npm view pnpm@10.1.0 dist.shasum) && corepack install -g pnpm@10.1.0+sha1.$SHA_SUM`;
  await $`pnpm i`;
}

export const runFrontendStage = (stage: "dev" | "prod") =>
  job.nothrow(`Frontend - ${stage}`, async () => {
    await docker.if(stage === "dev")(
      "node:22",
      { name: "Lint & Test" },
      async () => {
        await initPnpm();

        await $`pnpm -F web lint`;
        await $`pnpm -F web sync`; // TODO - we should run svelte-kit check here too
      },
    );

    await docker(
      "node:22",
      { name: "Build & Deploy" },
      async () => {
        await initPnpm();

        await $`pnpm -F web build`.env({
          PUBLIC_GITHUB_APP_NAME: stage === "dev"
            ? "pandaci-com-dev"
            : "pandaci-com",
          PUBLIC_API_URL: stage === "dev"
            ? `https://${env.PANDACI_BRANCH}.api.dev.pandaci.com`
            : "https://api.pandaci.com",
          PUBLIC_KRATOS_URL: stage === "dev"
            ? "https://auth.dev.pandaci.com"
            : "https://auth.pandaci.com",
          PUBLIC_STAGE: stage,
        });

        const projectName = stage === "dev"
          ? "pandaci-app-dev"
          : "pandaci-app-prod";

        const res =
          await $`pnpx wrangler pages deploy build --branch=${env.PANDACI_BRANCH} --commit-hash=${env.PANDACI_COMMIT_SHA} --project-name=${projectName}`
            .env({
              CLOUDFLARE_API_TOKEN: env.CLOUDFLARE_API_TOKEN,
              CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
            }).cwd("web").text();

        if (stage === "dev") {
          // We have a worker that rewrites our dev URLs to the correct deployment
          // we store the latest deployment for each sha in a KV store
          const deploymentId =
            /https:\/\/([a-zA-Z0-9-]+)\.pandaci-app-dev\.pages\.dev/.exec(res)
              ?.[1];

          await $`pnpx wrangler kv:key put ${env.PANDACI_BRANCH} ${deploymentId} --namespace-id=${env.CLOUDFLARE_KV_NAMESPACE_ID}`
            .env({
              CLOUDFLARE_API_TOKEN: env.CLOUDFLARE_API_TOKEN,
              CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
            });
        }
      },
    );
  });
