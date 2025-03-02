import type { Conclusion } from "@pandaci/proto";
import type { Volume } from "../volume.ts";

export interface DockerTaskOptions {
  /**
   * Override the derived name from the image
   */
  name?: string;

  /**
   * Docker volumes to be attached to the task
   * TODO - add example on how to create a volume
   */
  volumes?: Volume[];

  /**
   * Skip this task
   */
  skip?: boolean;

  /**
   *  If false, the task will not throw an error if the command fails,
   *  the task will still be marked as failed.
   *
   *  This is useful if you want to continue running the job even if a task fails
   *
   * @default true
   */
  throws?: boolean;
}

export interface DockerTaskResult {
  conclusion: Conclusion;
  isFailure: boolean;
  isSkipped: boolean;
  isSuccess: boolean;
  id: string;
  taskName: string;
}

export type DockerTaskFn = (() => void) | (() => Promise<void>);

type DockerTaskBase =
  & ((
    image: string,
    options: DockerTaskOptions,
    fn: DockerTaskFn,
  ) => Promise<DockerTaskResult>)
  & ((image: string, fn: DockerTaskFn) => Promise<DockerTaskResult>);

export type DockerMethod = "if" | "nothrow" | "skip";

type OmmitedDockerTask<T extends DockerMethod, K extends DockerMethod> =
  & Omit<
    DockerMethods<T | K>,
    T | K
  >
  & DockerTaskBase;

export interface DockerMethods<T extends DockerMethod = never>
  extends Record<DockerMethod, unknown> {
  if: (condition: boolean) => OmmitedDockerTask<T, "if">;
  nothrow: OmmitedDockerTask<T, "nothrow">;
  skip: OmmitedDockerTask<T, "skip">;
}

export type DockerTask = DockerTaskBase & DockerMethods;
