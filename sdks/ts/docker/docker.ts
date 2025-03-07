import { getWorkflowClient } from "../api.ts";
import { DockerTaskPromise } from "./dockerPromise.ts";
import type {
  DockerMethod,
  DockerMethods,
  DockerTask,
  DockerTaskFn,
  DockerTaskOptions,
} from "./dockerTypes.ts";

const getDocker = (
  skipChains: DockerMethod[],
  opts?: DockerTaskOptions,
): DockerTask => {
  const dockerFn = (
    image: string,
    arg2: DockerTaskFn | DockerTaskOptions,
    arg3?: DockerTaskFn,
  ) =>
    new DockerTaskPromise(
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

  const methods = {} as DockerMethods;

  if (!skipChains.includes("if")) {
    methods.if = function (condition: boolean) {
      return getDocker(["if", ...skipChains], {
        ...opts,
        skip: !condition,
      });
    };
  }

  if (!skipChains.includes("nothrow")) {
    methods.nothrow = getDocker(["nothrow", ...skipChains], {
      ...opts,
      throws: false,
    });
  }

  if (!skipChains.includes("skip")) {
    methods.skip = getDocker(["skip", ...skipChains], { ...opts, skip: true });
  }

  return Object.assign(
    dockerFn,
    methods,
  ) as DockerTask & DockerMethods;
};

/**
 * Create a new docker task
 * 
 * You can nest multiple tasks inside a job or run them without one. If using without a job,
 * a job will be created automatically using the default runner (ubuntu-4x).
 * 
 * @param {string} image - The docker image to use
 * @param {DockerTaskFn} fn - The function to run inside the docker container
 * @param {DockerTaskOptions} options - The options for the task
 * 
 * @returns {DockerTask} DockerTask (extends Promise<DockerTaskResult>) 
 * 
 * @method if - Skip this task if the condition is false
 * @method nothrow - Do not throw an error if the task fails
 * @method skip - Skip this task
 *
 * @example
 * ```ts
 * docker("node:22", () => {
 *   // Your steps here 
 * }); 
 * ```
 */
export const docker: DockerTask = getDocker([]);
