/**
 * Represents the result of an execution of a command.
 */
export class ExecResult {
  private _stdout: Uint8Array;
  private _stderr: Uint8Array;
  private _stdall: Uint8Array;
  private _exitCode: number;

  constructor({
    stdout,
    stderr,
    stdall,
    exitCode,
  }: {
    stdout: Uint8Array;
    stderr: Uint8Array;
    stdall: Uint8Array;
    exitCode: number;
  }) {
    this._stdout = stdout;
    this._stderr = stderr;
    this._stdall = stdall;
    this._exitCode = exitCode;
  }

  get stdout(): Uint8Array {
    return this._stdout;
  }

  get stderr(): Uint8Array {
    return this._stderr;
  }

  get exitCode(): number {
    return this._exitCode;
  }

  get stdall(): Uint8Array {
    return this._stdall;
  }

  /**
   * Returns the text of the stdall, which is the combination of stdout and stderr.
   * @param encoding The encoding to use. Default is "utf-8".
   * @returns The text of the stdall. 
   */
  text(encoding: string = "utf-8"): string {
    return new TextDecoder(encoding).decode(this._stdall);
  }

  /**
   * Returns the text of the stdall, which is the combination of stdout and stderr.
   * If you need to change the encoding, use the `text` method.
   * @returns The text of the stdall. 
   */
  toString(): string {
    return this.text();
  }

  /**
   * Returns the JSON object of stdall, which is the combination of stdout and stderr. 
   * @returns The JSON object of the stdall. 
   */
  json<T = unknown>(): T {
    return JSON.parse(this.text());
  }

  /**
   * Returns the lines of the stdall as an array, which is the combination of stdout and stderr.
   * We split by `\r?\n` to support both Unix and Windows line endings.
   * @param encoding The encoding to use. Default is "utf-8".
   * @returns The lines of the stdout.
   */
  lines(encoding: string = "utf-8"): string[] {
    return this.text(encoding).split(/\r?\n/);
  }
}
