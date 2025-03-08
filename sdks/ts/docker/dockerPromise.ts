import {
  Conclusion as ProtoConclusion,
  StartTaskData_Docker_DockerVolume_Type,
  type TaskMeta,
  WorkflowAlert_Type,
} from "@pandaci/proto";
import type {
  DockerTaskFn,
  DockerTaskOptions,
  DockerTaskResult,
} from "./dockerTypes.ts";
import { setImmediate } from "node:timers";
import type { getWorkflowClient } from "../api.ts";
import { logger } from "../logger.ts";
import { jobContext, taskContext } from "../context.ts";
import { extendedEnv } from "../env.ts";
import { JobPromise } from "../job/jobPromise.ts";
import type { Volume } from "../volume.ts";
import { type Conclusion, protoConclusionToConclusion } from "../types.ts";

/**
 * Error class representing a failed Docker task.
 * 
 * @extends Error
 * 
 * @property {Conclusion} conclusion - The conclusion status of the Docker task
 * @property {boolean} isFailure - Indicates if the task resulted in a failure
 * @property {string} id - Unique identifier for the task
 * @property {string} taskName - The name of the Docker task
 * @property {boolean} isSuccess - Indicates if the task was successful
 * @property {DockerTaskResult} data - The complete task result data
 */
export class DockerTaskError extends Error {
  conclusion: Conclusion;
  isFailure: boolean;
  id: string;
  taskName: string;
  isSuccess: boolean;

  constructor(public data: DockerTaskResult) {
    super(`Task ${data.taskName} failed`);
    this.conclusion = data.conclusion;
    this.isFailure = data.isFailure;
    this.id = data.id;
    this.taskName = data.taskName;
    this.isSuccess = data.isSuccess;
  }
}

interface DockerTaskFunctionContext {
  getWorkflowClient: () => ReturnType<typeof getWorkflowClient>;
}


/**
 * Represents a running docker task within a workflow.
 * 
 * DockerTaskPromise extends the native Promise class and provides additional functionality
 * 
 * @extends Promise<DockerTaskResult>
 * 
 */
export class DockerTaskPromise extends Promise<DockerTaskResult> {
  private reject: (reason?: DockerTaskError) => void;
  private resolve: (value: DockerTaskResult) => void;

  private name: string;
  private image: string;

  private volumes: Volume[];
  private throws: boolean;
  private skip: boolean;

  constructor(
    ctx: DockerTaskFunctionContext,
    image: string,
    fn: (() => void) | (() => Promise<void>),
    options?: DockerTaskOptions,
  ) {
    let tempReject: (reason?: DockerTaskError) => void = () => {};
    let tempResolve: (value: DockerTaskResult) => void = () => {};
    super((resolve, reject) => {
      tempReject = reject;
      tempResolve = resolve;

      // internally, .then() calls the executor function and requires it to resolve or reject.
      // We don't include this in our type definition because it'll just confuse users.
      if (typeof ctx === "function") (ctx as any)(resolve, reject);
    });

    this.reject = tempReject;
    this.resolve = tempResolve;

    this.name = options?.name || image;
    this.throws = options?.throws ?? true;
    this.skip = options?.skip ?? false;
    this.image = image;
    this.volumes = options?.volumes ?? [];

    // If this is an internal .then() call then calling run will just
    if (typeof ctx === "function") return;

    // TODO - this can cause some issues since not all errors are DockerTaskErrors
    setImmediate(() => this.run(ctx, fn).catch(this.reject));
  }

  private async run(
    ctx: DockerTaskFunctionContext,
    fn: DockerTaskFn,
  ) {
    logger.info(`Starting docker task ${this.name}`);

    if (taskContext.getStore()) {
      await ctx.getWorkflowClient().createWorkflowAlert({
        workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
        alert: {
          type: WorkflowAlert_Type.ERROR,
          title: "Syntax error - Nested task detected",
          message: "Tasks cannot be nested inside other tasks",
        },
      });

      throw new Error("nested task");
    }

    const jobCtx = jobContext.getStore();

    if (!jobCtx) {
      // If we're not in a job, then we must create one
      await new JobPromise(
        ctx,
        this.name,
        () =>
          new DockerTaskPromise(ctx, this.image, fn, {
            name: this.name,
            skip: this.skip,
            throws: this.throws,
            volumes: this.volumes,
          })
            .then((res) => this.resolve(res))
            .catch((err) => this.reject(err)),
        {
          skip: this.skip,
          throws: false,
        },
      );
      return;
    }

    const { meta: jobMeta, registerJobPromise } = jobCtx;

    registerJobPromise(this);

    const workflowClient = ctx.getWorkflowClient();

    const startTaskResult = await workflowClient.startTask({
      $typeName: "proto.v1.WorkflowServiceStartTaskRequest",
      jobMeta,
      skipped: this.skip,
      workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
      data: {
        $typeName: "proto.v1.StartTaskData",
        name: this.name,
        data: {
          case: "dockerData",
          value: {
            $typeName: "proto.v1.StartTaskData.Docker",
            image: this.image,
            volumes: this.volumes.map(({ source, target, type }) => ({
              $typeName: "proto.v1.StartTaskData.Docker.DockerVolume",
              source,
              target,
              type: type === "bind"
                ? StartTaskData_Docker_DockerVolume_Type.BIND
                : StartTaskData_Docker_DockerVolume_Type.VOLUME,
            })) ?? [],
          },
        },
      },
    });

    if (this.skip) {
      logger.info(`Task ${this.name} skipped`);
      return this.resolve({
        conclusion: protoConclusionToConclusion(ProtoConclusion.SKIPPED),
        isFailure: false,
        id: startTaskResult.taskMeta!.id,
        taskName: this.name,
        isSkipped: true,
        isSuccess: false,
      });
    }

    if (startTaskResult.taskMeta?.specificMeta.case !== "dockerMeta") {
      await workflowClient.createWorkflowAlert({
        workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
        alert: {
          type: WorkflowAlert_Type.ERROR,
          title: "Syntax error - Wrong task type",
          message:
            `Wrong task type, expected dockerData but got ${startTaskResult.taskMeta?.specificMeta.case} instead`,
        },
      });

      throw new Error(
        `Wrong task type, expected dockerData but got ${startTaskResult.taskMeta?.specificMeta.case} instead`,
      );
    }

    const taskMeta: TaskMeta = {
      $typeName: "proto.v1.TaskMeta",
      id: startTaskResult.taskMeta!.id,
      name: this.name,
      specificMeta: {
        case: "dockerMeta",
        value: {
          $typeName: "proto.v1.TaskMeta.Docker",
          containerId: startTaskResult.taskMeta.specificMeta.value.containerId,
        },
      },
    };

    const promises: Promise<unknown>[] = [];

    // Run our task
    const conclusion = await taskContext
      .run(
        {
          meta: taskMeta,
          registerTaskPromise: (promise) => {
            promises.push(promise);
          },
        },
        async () => {
          try {
            await Promise.resolve(fn());
            const res = await Promise.allSettled(promises);
            if (res.some((r) => r.status === "rejected")) {
              logger.error(`Task ${this.name} failed`);
              return ProtoConclusion.FAILURE;
            }
          } catch (err) {
            // TODO - we should check the error and throw if it's not due to a non-zero exit code
            logger.error(`Task ${this.name} failed`, err);
            return ProtoConclusion.FAILURE;
          }
          return ProtoConclusion.SUCCESS;
        },
      );

    // There's not much we can do if this throws an error since
    // it we'd need to call that endpoint to handle it anyway
    await workflowClient.stopTask({
      $typeName: "proto.v1.WorkflowServiceStopTaskRequest",
      workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
      jobMeta,
      conclusion,
      taskMeta,
    });

    logger.info(`Job ${this.name} completed`);

    const result: DockerTaskResult = {
      conclusion: protoConclusionToConclusion(conclusion),
      isFailure: conclusion === ProtoConclusion.FAILURE,
      id: taskMeta.id,
      taskName: this.name,
      isSkipped: false,
      isSuccess: conclusion === ProtoConclusion.SUCCESS,
    };

    if (this.throws && conclusion === ProtoConclusion.FAILURE) {
      this.reject(new DockerTaskError(result));
    }

    this.resolve(result);
  }
}
