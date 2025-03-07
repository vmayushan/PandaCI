import { WorkflowAlert_Type, Conclusion as ProtoConclusion } from "@pandaci/proto";
import type { JobOptions, JobResult } from "./jobTypes.ts";
import { setImmediate } from "node:timers";
import { jobContext } from "../context.ts";
import { extendedEnv } from "../env.ts";
import { logger } from "../logger.ts";
import type { JobFn } from "./mod.ts";
import type { getWorkflowClient } from "../api.ts";
import { type Conclusion, protoConclusionToConclusion } from "../types.ts";

/**
 * Represents an error that occurred during job execution.
 * Extends the standard Error class with additional job-related properties.
 * 
 * @extends Error
 * 
 * @property {Conclusion} conclusion - The conclusion status of the job
 * @property {boolean} isFailure - Indicates if the job failed
 * @property {boolean} isSkipped - Indicates if the job was skipped
 * @property {boolean} isSuccess - Indicates if the job completed successfully
 * @property {string} id - The unique identifier of the job
 * @property {string} jobName - The name of the job
 * @property {string} [runner] - The runner that executed the job, if available
 * 
 * @param {JobResult} result - The result object containing job execution details
 */
export class JobError extends Error {
  conclusion: Conclusion;
  isFailure: boolean;
  isSkipped: boolean;
  isSuccess: boolean;
  id: string;
  jobName: string;
  runner?: string;

  constructor(public result: JobResult) {
    super(`Job ${result.name} failed`);
    this.conclusion = result.conclusion;
    this.isFailure = result.isFailure;
    this.id = result.id;
    this.jobName = result.name;
    this.runner = result.runner;
    this.isSkipped = result.isSkipped;
    this.isSuccess = result.isSuccess;
  }
}

interface JobFunctionContext {
  getWorkflowClient: () => ReturnType<typeof getWorkflowClient>;
}


/**
 * A Promise-based representation of a job that can be executed in a workflow.
 * 
 * JobPromise extends the native Promise class and provides additional functionality
 * for tracking and controlling job execution within a workflow context. It handles
 * job creation, execution, and completion while maintaining the standard Promise interface.
 * 
 * @extends Promise<JobResult>
 * 
 */
export class JobPromise extends Promise<JobResult> {
  private reject: (reason?: JobResult) => void;
  private resolve: (value: JobResult) => void;

  private name: string;
  private runner: string;

  private throws?: boolean;
  private skip?: boolean;

  constructor(
    ctx: JobFunctionContext,
    name: string,
    fn: (() => void) | (() => Promise<void>),
    options?: JobOptions,
  ) {
    let tempReject: (reason?: JobResult) => void = () => {};
    let tempResolve: (value: JobResult) => void = () => {};
    super((resolve, reject) => {
      tempReject = reject;
      tempResolve = resolve;

      // internally, .then() calls the executor function and requires it to resolve or reject.
      // We don't include this in our type definition because it'll just confuse users.
      if (typeof ctx === "function") (ctx as any)(resolve, reject);
    });

    this.reject = tempReject;
    this.resolve = tempResolve;

    this.name = name;
    this.runner = options?.runner ?? "ubuntu-4x";
    this.throws = options?.throws ?? true;
    this.skip = options?.skip ?? false;

    // If this is an internal .then() call then calling run will just
    if (typeof ctx === "function") return;

    // TODO - this can cause some issues since not all errors are JobError
    setImmediate(() => this.run(ctx, fn).catch(this.reject));
  }

  private async run(
    ctx: JobFunctionContext,
    fn: JobFn,
  ) {
    logger.info(`Starting job ${this.name}`);

    if (jobContext.getStore()) {
      await ctx.getWorkflowClient().createWorkflowAlert({
        workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
        alert: {
          type: WorkflowAlert_Type.ERROR,
          title: "Syntax error - Nested job detected",
          message: "Jobs cannot be nested inside other jobs or tasks",
        },
      });

      throw new Error("nested job");
    }

    const client = ctx.getWorkflowClient();

    const createJobResult = await client.startJob({
      name: this.name,
      runner: this.runner,
      workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
      skipped: this.skip,
    });

    if (this.skip) {
      logger.info(`Job ${this.name} skipped`);
      return this.resolve({
        conclusion: "skipped",
        id: createJobResult.jobMeta!.id,
        name: createJobResult.jobMeta!.name,
        isFailure: false,
        runner: createJobResult.jobMeta!.runner,
        isSkipped: true,
        isSuccess: false,
      });
    }

    const taskPromises: Promise<unknown>[] = [];

    const conclusion = await jobContext
      .run(
        {
          meta: createJobResult.jobMeta!,
          registerJobPromise: (promise) => {
            taskPromises.push(promise);
          },
        },
        async () => {
          try {
            await Promise.resolve(fn());
            await Promise.allSettled(taskPromises);
          } catch (err) {
            logger.error(`Job ${this.name} failed`, err);
            return ProtoConclusion.FAILURE;
          }
          return ProtoConclusion.UNSPECIFIED;
        },
      );

    const jobRes = await client.stopJob({
      $typeName: "proto.v1.WorkflowServiceStopJobRequest",
      jobMeta: createJobResult.jobMeta!,
      workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
      conclusion,
    });

    logger.info(`Job ${this.name} completed`);

    const res: JobResult = {
      conclusion: protoConclusionToConclusion(jobRes.conclusion),
      id: createJobResult.jobMeta!.id,
      name: this.name,
      runner: createJobResult.jobMeta!.runner,
      isFailure: jobRes.conclusion === ProtoConclusion.FAILURE,
      isSkipped: jobRes.conclusion === ProtoConclusion.SKIPPED,
      isSuccess: jobRes.conclusion === ProtoConclusion.SUCCESS,
    };

    if (this.throws && res.isFailure) {
      this.reject(new JobError(res));
    }

    this.resolve(res);
  }
}
