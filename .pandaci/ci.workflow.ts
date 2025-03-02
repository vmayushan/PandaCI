import { Config, env } from "@pandaci/workflow";
import { runBackendStage } from "./backend/stage.ts";
import { runFrontendStage } from "./frontend/stage.ts";

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

const devJobs = await Promise.all([
  runFrontendStage("dev"),
  runBackendStage("dev"),
]);

if (env.PANDACI_BRANCH === "main" && devJobs.every((job) => !job.isFailure)) {
  await Promise.all([
    runFrontendStage("prod"),
    runBackendStage("prod"),
  ]);
}
