import { getWorkflowClient } from "./api.ts";
import { jobContext } from "./context.ts";
import { extendedEnv } from "./env.ts";
import { logger } from "./logger.ts";

export interface Volume {
  source: string;
  target: string;
  type: "bind" | "volume";
}

export async function volume({ host }: { host?: string } = {}): Promise<
  (path: string) => Volume
> {
  const jobCtx = jobContext.getStore();

  if (!jobCtx) {
    logger.error("Task must be started within a job");
    throw new Error("Task must be started within a job");
  }

  const { meta: jobMeta } = jobCtx;

  const volumeRes = await getWorkflowClient().createJobVolume({
    $typeName: "proto.v1.WorkflowServiceCreateJobVolumeRequest",
    host,
    jobMeta,
    workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
  });

  return (target: string) => ({
    source: volumeRes.source,
    target,
    type: host ? "bind" : "volume",
  });
}
