import type { Conclusion } from "@pandaci/proto";
import type { JobPromise } from "./jobPromise.ts";

export type JobFn = (() => void) | (() => Promise<void>);

export type Runner =
  | "ubuntu-1x"
  | "ubuntu-2x"
  | "ubuntu-4x"
  | "ubuntu-8x"
  | "ubuntu-16x";

export interface JobOptions {
  /**
   * The Runner to use for this job.
   * @default "ubuntu-4x"
   */
  runner?: Runner;

  /**
   * Skip this job
   */
  skip?: boolean;

  /**
   * If the job should throw an error if it fails
   * @default true
   */
  throws?: boolean;
}

export interface JobResult {
  conclusion: Conclusion;
  isFailure: boolean;
  isSuccess: boolean;
  isSkipped: boolean;
  id: string;
  name: string;
  runner?: string;
}

type JobBase =
  & ((
    name: string,
    options: JobOptions,
    fn: JobFn,
  ) => JobPromise)
  & ((image: string, fn: JobFn) => JobPromise);

export type JobMethod = "if" | "nothrow" | "skip";

type OmmitedJob<T extends JobMethod, K extends JobMethod> =
  & Omit<
    JobMethods<T | K>,
    T | K
  >
  & JobBase;

export interface JobMethods<T extends JobMethod = never>
  extends Record<JobMethod, unknown> {
  if: (condition: boolean) => OmmitedJob<T, "if">;
  nothrow: OmmitedJob<T, "nothrow">;
  skip: OmmitedJob<T, "skip">;
}

export type Job = JobBase & JobMethods;
