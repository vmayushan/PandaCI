import { AsyncLocalStorage } from "node:async_hooks";
import { buildCmd } from "./command.ts";
import { type ExecOptions, ExecPromise } from "./execPromise.ts";
import {
  type JobContext,
  jobContext,
  type TaskContext,
  taskContext,
} from "../context.ts";
import { extendedEnv } from "../env.ts";
import { getWorkflowClient } from "../api.ts";

export function isOptionsObject(options: unknown): boolean {
  return (
    Boolean(options) && typeof options === "object" && !Array.isArray(options)
  );
}

export interface Exec {
  (pieces: TemplateStringsArray, ...args: unknown[]): ExecPromise;
  (opts: ExecOptions): Exec;
}

export interface ShellContext {
  options: ExecOptions;
  jobCtx?: JobContext;
  taskCtx?: TaskContext;
}

export const shellContext = new AsyncLocalStorage<ShellContext>();

export const $: Exec = ((
  firstArg: TemplateStringsArray | ExecOptions,
  ...args: unknown[]
) => {
  const initialJobCtx = jobContext.getStore();
  const initalTaskCtx = taskContext.getStore();

  // TODO - we should probably be able to avoid having this weird inital context stuff
  const { options, jobCtx, taskCtx } = shellContext.getStore() || {
    options: {},
    jobCtx: initialJobCtx,
    taskCtx: initalTaskCtx,
  };

  if (isOptionsObject(firstArg)) {
    Object.assign(options, firstArg);

    return (...args: unknown[]) =>
      shellContext.run(
        { jobCtx, options, taskCtx },
        () => $(...(args as [TemplateStringsArray, ...unknown[]])),
      );
  }

  const cmd = buildCmd(firstArg as TemplateStringsArray, args);

  const promise = new ExecPromise(
    {
      jobContext: jobCtx,
      taskContext: taskCtx,
      workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
      client: getWorkflowClient(),
    },
    cmd,
    options,
  );

  taskCtx?.registerTaskPromise(promise);

  return promise;
}) as Exec;
