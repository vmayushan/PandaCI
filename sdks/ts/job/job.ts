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

export const job: Job = getJob([]);
