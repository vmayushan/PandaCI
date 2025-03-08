import { getWorkflowClient } from "./api.ts";
import { jobContext } from "./context.ts";
import { extendedEnv } from "./env.ts";
import { logger } from "./logger.ts";

/**
 * Docker volume 
 */
export interface Volume {
  /**
   * Source path
   * @example /path/to/source
   */
  source: string;
  /**
   * Target path
   * @example /path/to/target
   */
  target: string;
  /**
   * Type of volume
   * 
   * - `bind` - bind mount
   * - `volume` - docker volume
   */
  type: "bind" | "volume";
}

/**
 * Create a volume. Must be used within a job. Volumes cannot be created outside of a job or used across jobs.
 * @param host Host path
 * @returns A function that creates a volume
 */
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
