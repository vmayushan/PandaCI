export * from "./docker/mod.ts";
export * from "./job/mod.ts";
export { type Env, env } from "./env.ts";
export * from "./shell/mod.ts";
export { type Volume, volume } from "./volume.ts";
export type { ExecOptions, ExecResult } from "./shell/mod.ts";
export type { Conclusion } from "./types.ts";

/**
 * Pull request event type
 */
export type PullRequestEvent = "opened" | "closed" | "reopened" | "synchronize";

/**
 * Workflow config
 * @example export const config: Config = { name: "My Job", on: { push: { branches: ["main"] } } }
 */
export interface Config {
  /**
   * Job name, defaults to the filename and path.
   */
  name?: string;
  /**
   * Trigger configuration, defaults to running on every push.
   *
   * @default { push: { branches: ["*"] } }
   */
  on?: {
    /**
     * Push configuration, by default all branches are triggered.
     * If you want to trigger on specific branches, you need to specify the branches.
     */
    push?: {
      /**
       * Branches to trigger on, defaults to all branches.
       * You can use glob patterns to match multiple branches.
       * @example ["main", "feature/*"]
       */
      branches?: string[];
      /**
       * Branches to ignore, defaults to none.
       * You can use glob patterns to match multiple branches.
       * @example ["main", "feature/*"]
       */
      branchesIgnore?: string[];
    };
    /**
     * Pull request configuration, by default no pull requests are triggered.
     * If you want to trigger on pull requests, you need to specify the events.
     */
    pr?: {
      /**
       * Events to trigger on, defaults to none.
       * @example ["opened", "closed"]
       */
      events?: PullRequestEvent[];
      /**
       * Branches to trigger on, defaults to all branches.
       * You can use glob patterns to match multiple branches.
       * @example ["main", "feature/*"]
       */
      targetBranches?: string[];
      /**
       * Branches to ignore, defaults to none.
       * You can use glob patterns to match multiple branches.
       * @example ["main", "feature/*"]
       */
      targetBranchesIgnore?: string[];
    };
  };
}

// Caching behavior:
// 1. Auto save cache on successfull job completion (probably offer a way to disable this)
// 2. Expose a way to manually save cache
// 3. Expose a way to manually restore cache



// export function Cache({}: {
//   keys: string[];
//   paths: string[];
// }) {
//   return Promise.resolve();
// }

export class Cache {



  constructor({
    keys,
    paths,
    autoSave = true,
    autoRestore = true,
  }: {
    // Cache keys, we go over each key in order until we find a match
    keys: string[];
    // Paths to cache, supports glob patterns
    paths: string[]; 
    // Automatically save cache on successfull job completion, defaults to true
    autoSave?: boolean;
    // Automatically restore cache on job start, defaults to true
    autoRestore?: boolean;
  }) {}

  async restore() {
    // Check if cache exists
    // If cache exists, restore it
    // TODO - maybe display a warning if we have already restored a cache / are overwriting files
  }

  async save() {
    // Save cache
  }

  async clear() {
    // TODO - decide if we want to have this method
  }

  get hit() {
    // Return true if cache was restored
    return false;
  }

  get miss() {
    // opposite of hit
    return !this.hit;
  }

  get key() {
    // Return the key that was used to restore the cache
    return "";
  } 


}

const npmCache = new Cache({
  keys: [""],
  paths: ["node_modules"],
})

