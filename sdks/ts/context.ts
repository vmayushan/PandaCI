import { AsyncLocalStorage } from "node:async_hooks";
import type { JobMeta, TaskMeta } from "@pandaci/proto";

export interface JobContext {
  meta: JobMeta;
  registerJobPromise: (promise: Promise<unknown>) => void;
}

export const jobContext = new AsyncLocalStorage<JobContext>();

export interface TaskContext {
  meta: TaskMeta;
  registerTaskPromise: (promise: Promise<unknown>) => void;
}

export const taskContext = new AsyncLocalStorage<TaskContext>();
