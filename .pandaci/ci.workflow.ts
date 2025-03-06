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

await Promise.all([
  job("Backend dev", async () => {
    await testGoTask();
    await deployDevCoreBackendTask();
  }),
  job("Frontend dev", async () => {
    await testFrontendTask();
    await deployDevFrontendTask();
  }),
]);

await Promise.all([
  deployProdCoreBackendTask(),
  deployProdFrontendTask(),
]);
