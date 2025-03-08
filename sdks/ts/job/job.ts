import { getWorkflowClient } from "../api.ts";
import { JobPromise } from "./jobPromise.ts";
import type {
  Job,
  JobFn,
  JobMethod,
  JobMethods,
  JobOptions,
} from "./jobTypes.ts";

const getJob = (
  skipChains: JobMethod[],
  opts?: JobOptions,
): Job => {
  const jobFn = (image: string, arg2: JobFn | JobOptions, arg3?: JobFn) =>
    new JobPromise(
      {
        getWorkflowClient,
      },
      image,
      typeof arg2 === "function" ? arg2 : arg3!,
      {
        ...opts,
        ...(typeof arg2 === "function" ? {} : arg2),
      },
    );

  const methods = {} as JobMethods;

  if (!skipChains.includes("if")) {
    methods.if = function (condition: boolean) {
      return getJob(["if", ...skipChains], {
        ...opts,
        skip: !condition,
      });
    };
  }

  if (!skipChains.includes("nothrow")) {
    methods.nothrow = getJob(["nothrow", ...skipChains], {
      ...opts,
      throws: false,
    });
  }

  if (!skipChains.includes("skip")) {
    methods.skip = getJob(["skip", ...skipChains], { ...opts, skip: true });
  }

  return Object.assign(
    jobFn,
    methods,
  ) as Job & JobMethods;
};

/**
 * Create a new job
 *
 * @returns {JobPromise} JobPromise (extends Promise<JobResult>)
 *
 * @method if - Skip this job if the condition is false
 * @method nothrow - Do not throw an error if this job fails
 * @method skip - Skip this job
 *
 * @example ```
 * job("my-job", () => {
 *    // your job code here
 * });
 * ```
 *
 * @example ```
 * // With options
 * job.nothrow("my-job", { runner: 'ubuntu-2x' }, () => {
 *    // your job code here
 * });
 * ```
 */
export const job: Job = getJob([]);
