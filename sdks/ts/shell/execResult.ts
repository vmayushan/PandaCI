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

  text(encoding: string = "utf-8"): string {
    return new TextDecoder(encoding).decode(this._stdall);
  }

  toString(): string {
    return this.text();
  }

  json<T = unknown>(): T {
    return JSON.parse(this.text());
  }

  lines(encoding: string = "utf-8"): string[] {
    return this.text(encoding).split(/\r?\n/);
  }
}
