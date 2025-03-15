import type { Conclusion } from "../types.ts";
import type { JobPromise } from "./jobPromise.ts";

/**
* A function that can be executed as a job
* It can be a synchronous function that returns void
* Or an asynchronous function that returns Promise<void>
*/
export type JobFn = (() => void) | (() => Promise<void>);

/**
 * The runner to use for a job
 */
export type Runner =
  | "ubuntu-1x"
  | "ubuntu-2x"
  | "ubuntu-4x"
  | "ubuntu-8x"
  | "ubuntu-16x";

/**
 * Options for a job
 */
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

/**
 * The result of a job. This is what is returned when a job is finished executing
 *
 */
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

/**
 * The methods available on a job
*/
export interface JobMethods<T extends JobMethod = never>
  extends Record<JobMethod, unknown> {
  if: (condition: boolean) => OmmitedJob<T, "if">;
  nothrow: OmmitedJob<T, "nothrow">;
  skip: OmmitedJob<T, "skip">;
}

/**
 * Create a new job
 * @returns {JobPromise} JobPromise (extends Promise<JobResult>)
 * @method if - Skip this job if the condition is false
 * @method nothrow - Do not throw an error if this job fails
 * @method skip - Skip this job
 *
 * @example
 * ```
 * job("my-job", () => {
 * // ...
 * });
 */
export type Job = JobBase & JobMethods;
