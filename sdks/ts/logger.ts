import { extendedEnv } from "./env.ts";

export interface Logger {
  panic: (message: string) => void;
  fatal: (message: string) => void;
  error: (message: string, err?: unknown) => void;
  warn: (message: string) => void;
  info: (message: string) => void;
  debug: (message: string) => void;
  trace: (message: string) => void;
}

const logLevels = {
  panic: 5,
  fatal: 4,
  error: 3,
  warn: 2,
  info: 1,
  debug: 0,
  trace: -1,
};

const logLevel = extendedEnv.PANDACI_LOG_LEVEL ?? -1;

export const logger: Logger = {
  panic: (message: string) => {
    if (logLevels.panic >= logLevel) console.error(`PANIC: ${message}`);
  },
  fatal: (message: string) => {
    if (logLevels.fatal >= logLevel) console.error(`FATAL: ${message}`);
  },
  error: (message: string, err?: unknown) => {
    if (logLevels.error >= logLevel) {
      console.error(`ERROR: ${message}`, ...(err ? [err] : []));
    }
  },
  warn: (message: string) => {
    if (logLevels.warn >= logLevel) console.warn(`WARN: ${message}`);
  },
  info: (message: string) => {
    if (logLevels.info >= logLevel) console.log(`INFO: ${message}`);
  },
  debug: (message: string) => {
    if (logLevels.debug >= logLevel) console.log(`DEBUG: ${message}`);
  },
  trace: (message: string) => {
    if (logLevels.trace >= logLevel) console.log(`TRACE: ${message}`);
  },
};
