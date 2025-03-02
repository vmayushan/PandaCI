import { $, Config, docker, env } from "@pandaci/workflow";

export const config: Config = {
  name: "Cleanup",
  on: {
    push: {
      branches: [],
    },
    pr: {
      events: ["closed"],
      targetBranches: ["main"],
    },
  },
};

await docker("ubuntu:24.04", { name: "Cleanup fly" }, async () => {
  await $`apt-get update && apt-get install -y curl`;

  await $`curl -L https://fly.io/install.sh | sh`;

  // make fly cli available from 'flyctl' command
  await $`ln -s /root/.fly/bin/flyctl /usr/local/bin/flyctl`;

  const flyEnv = {
    FLY_ACCESS_TOKEN: env.FLY_ACCESS_TOKEN_DEV,
    FLY_APP: `pandaci-core-prev-${env.PANDACI_BRANCH}-${env.PANDACI_PR_NUMBER}`,
  };

  if (
    (await $`flyctl status`.env(flyEnv)
      .nothrow()).exitCode === 0
  ) {
    await $`flyctl apps delete ${flyEnv.FLY_APP} --yes`.env(flyEnv);
  }

  // TODO - cleanup frontend
});
