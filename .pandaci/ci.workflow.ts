// Note: This is normally "jsr:@pandaci/workflow", but we're just using the local source here.
import { Config, job } from "@pandaci/workflow";
import {
  deployDevCoreBackendTask,
  deployProdCoreBackendTask,
  testGoTask,
} from "./ci/backend.ts";
import {
  deployDevFrontendTask,
  deployProdFrontendTask,
  testFrontendTask,
} from "./ci/frontend.ts";
import { env } from "./utils.ts";

export const config: Config = {
  name: "CI",
  on: {
    push: {
      branches: ["main"],
    },
    pr: {
      events: ["synchronize", "opened", "reopened"],
      targetBranches: ["main"],
    },
  },
};

const devRes = await Promise.allSettled([
  job("Backend dev", async () => {
    await testGoTask();
    await deployDevCoreBackendTask();
  }),
  job("Frontend dev", async () => {
    await testFrontendTask();
    await deployDevFrontendTask();
  }),
]);

if (devRes.some((result) => result.status === "rejected")) {
  throw new Error("One or more jobs failed");
}

if (env.PANDACI_BRANCH === "main") {
  const prodRes = await Promise.allSettled([
    deployProdCoreBackendTask(),
    deployProdFrontendTask(),
  ]);

  if (prodRes.some((result) => result.status === "rejected")) {
    throw new Error("One or more jobs failed");
  }
}
