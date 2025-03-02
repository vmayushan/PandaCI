// Original source:
// https://github.com/google/zx/blob/main/src/util.ts

export function quote(arg: string): string {
  if (/^[\w/.\-@:=]+$/.test(arg) || arg === "") {
    return arg;
  }
  return (
    `$'` +
    arg
      .replace(/\\/g, "\\\\")
      .replace(/'/g, "\\'")
      .replace(/\f/g, "\\f")
      .replace(/\n/g, "\\n")
      .replace(/\r/g, "\\r")
      .replace(/\t/g, "\\t")
      .replace(/\v/g, "\\v")
      .replace(/\0/g, "\\0") +
    `'`
  );
}

export function quotePowerShell(arg: string): string {
  if (/^[\w/.\-]+$/.test(arg) || arg === "") {
    return arg;
  }
  return `'` + arg.replace(/'/g, "''") + `'`;
}
