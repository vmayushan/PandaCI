import { quote } from "./quote.ts";

// deno-lint-ignore no-explicit-any
const isPromiseLike = (value: any): boolean =>
  typeof value?.then === "function";

type Arg = { stdout?: unknown } & object;

const substitute = (arg: Arg) =>
  arg?.stdout instanceof Uint8Array
    ? new TextDecoder().decode(arg.stdout)
    : `${arg}`;

export type EscapeFn = (input: string) => string;
export const buildCmd = (
  pieces: TemplateStringsArray,
  args: unknown[],
): string | Promise<string> => {
  // If the user is piping in a promise, wait for it to resolve
  if (args.some(isPromiseLike)) {
    return Promise.all(args).then((args) => buildCmd(pieces, args));
  }

  let cmd = pieces[0] || "",
    i = 0;
  while (i < args.length) {
    const s = Array.isArray(args[i])
      ? (args[i] as unknown[])
        .map((x) => quote(substitute(x as Arg)))
        .join(" ")
      : quote(substitute(args[i] as Arg));

    cmd += s + pieces[++i];
  }

  return cmd;
};
