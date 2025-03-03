import { WorkflowAlert_Type } from "@pandaci/proto";
import type { JobContext, TaskContext } from "../context.ts";
import { setImmediate } from "node:timers";
import { getWorkflowClient } from "../api.ts";
import { extendedEnv } from "../env.ts";
import { ExecResult } from "./execResult.ts";
import { logger } from "../logger.ts";

function mergeUint8Array(
  ...arrays: Uint8Array[]
): Uint8Array<ArrayBuffer> {
  const length = arrays.reduce((acc, arr) => acc + arr.length, 0);
  const result = new Uint8Array(length);

  let offset = 0;
  for (const arr of arrays) {
    result.set(arr, offset);
    offset += arr.length;
  }

  return result;
}

export interface ExecOptions {
  nothrow?: boolean;
  cwd?: string;
  env?: Record<string, string | number>;
}

export interface ExecPromiseContext {
  workflowJwt: string;
  jobContext?: JobContext;
  taskContext?: TaskContext;
  client: ReturnType<typeof getWorkflowClient>;
}

export class ExecPromise extends Promise<ExecResult> {
  private _reject: (reason?: ExecResult) => void;
  private _resolve: (value: ExecResult) => void;

  private _cmd: string | Promise<string>;

  private _nothrow?: boolean;
  private _cwd?: string;
  private _env?: Record<string, string | number>;

  constructor(
    firstArg: ExecPromiseContext,
    cmd: string | Promise<string>,
    options?: ExecOptions,
  ) {
    let tempReject: (reason?: ExecResult) => void = () => {};
    let tempResolve: (value: ExecResult) => void = () => {};
    super((resolve, reject) => {
      tempReject = reject;
      tempResolve = resolve;

      // internally, .then() calls the executor function and requires it to resolve or reject.
      // We don't include this in our type definition because it'll just confuse users.
      if (typeof firstArg === "function") (firstArg as any)(resolve, reject);
    });

    this._reject = tempReject;
    this._resolve = tempResolve;

    this._cmd = cmd;
    this._nothrow = options?.nothrow;
    this._cwd = options?.cwd;
    this._env = options?.env;

    // If this is an internal .then() call then calling run will just
    if (typeof firstArg === "function") return;

    setImmediate(() => this.run(firstArg).catch(this._reject));
  }

  private async run(
    { workflowJwt, jobContext, taskContext, client }: ExecPromiseContext,
  ) {
    let stdout = new Uint8Array();
    let stderr = new Uint8Array();
    let stdall = new Uint8Array();
    let exitCode = 0;

    if (!taskContext || !jobContext) {
      await getWorkflowClient().createWorkflowAlert({
        workflowJwt: extendedEnv.PANDACI_WORKFLOW_JWT,
        alert: {
          type: WorkflowAlert_Type.ERROR,
          title: "Syntax error - Exec detected outside of task",
          message: "Exec must be called inside a task context",
        },
      });
      throw new Error("Exec called outside of task context");
    }

    const stdoutChunks: Uint8Array[] = [];
    const stderrChunks: Uint8Array[] = [];
    const stdallChunks: Uint8Array[] = [];

    // The output gets sent in chunks, so we need to wait for the final exit code
    for await (
      const res of client
        .startStep(
          {
            $typeName: "proto.v1.WorkflowServiceStartStepRequest",
            jobMeta: jobContext.meta,
            workflowJwt,
            taskMeta: taskContext.meta,
            data: {
              case: "execData",
              value: {
                $typeName: "proto.v1.ExecInfo",
                cmd: await this._cmd,
                cwd: this._cwd || "",
                // TODO - add a warning if a value is undefined, the user might have forgotten to set it
                // if possible, we should do this at the env. level since we might want to set an undefined value to override a previous value
                env: Object.entries(this._env || {}).filter(([_, value]) =>
                  value !== undefined
                ).map(([key, value]) => ({
                  key,
                  value: value.toString(),
                  $typeName: "proto.v1.EnvironmentVariable",
                })),
              },
            },
          },
        )
    ) {
      if (res.payload?.case !== "exec") {
        throw new Error("Wrong step result type. exec expected");
      }

      if (res.payload.value.logData.case === "stdout") {
        stdoutChunks.push(res.payload.value.logData.value);
        stdallChunks.push(res.payload.value.logData.value);
      } else if (res.payload.value.logData.case === "stderr") {
        stderrChunks.push(res.payload.value.logData.value);
        stdallChunks.push(res.payload.value.logData.value);
      } else if (res.payload.value.logData.case === "exitCode") {
        exitCode = res.payload.value.logData.value;
      }
    }

    stdout = mergeUint8Array(...stdoutChunks);
    stderr = mergeUint8Array(...stderrChunks);
    stdall = mergeUint8Array(...stdallChunks);

    const result = new ExecResult({
      stdout,
      stderr,
      stdall,
      exitCode,
    });

    if (exitCode !== 0 && !this._nothrow) {
      logger.error(result.text());
      this._reject(result);
    } else {
      this._resolve(result);
    }
  }

  nothrow(): this {
    this._nothrow = true;
    return this;
  }

  throws(throws: boolean = true): this {
    this._nothrow = !throws;
    return this;
  }

  cwd(cwd: string): this {
    this._cwd = cwd;
    return this;
  }

  env(env?: Record<string, string | number>): this {
    this._env = env;
    return this;
  }

  async text(encoding: string = "utf-8"): Promise<string> {
    const result = await this;
    return result.text(encoding);
  }

  async lines(encoding: string = "utf-8"): Promise<string[]> {
    const result = await this;
    return result.lines(encoding);
  }

  async json<T = unknown>(): Promise<T> {
    const result = await this;
    return result.json<T>();
  }
}
