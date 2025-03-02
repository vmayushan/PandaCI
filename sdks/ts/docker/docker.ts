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

export const docker: DockerTask = getDocker([]);
