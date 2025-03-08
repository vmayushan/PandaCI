# @pandaci/core


Checkout https://pandaci.com/docs/typescript-syntax/overview/overview/ for full docs

```typescript
import { docker, $, job, env } from "jsr:@pandaci/workflow";

// Creates a docker task on a 4-core machine
// same as wrapping with `job("some-name", () => { ... })`
docker("ubuntu:latest", { name: "hello world" }, () => {
  // Runs a shell command
  $`echo "Hello, world! from branch: ${env.PANDACI_BRANCH}"`;
});

// Creates 2 docker tasks on an 2-core machine
job("Job 2", { runner: "ubuntu-2x" }, async () => {
  let files = '';
  await docker("ubuntu:latest", async () => {
    // Runs the command in the .pandaci directory
    // and stores the output as a string
    files = await $`ls`.cwd(".pandaci").text();
  });
  docker("ubuntu:latest", () => {
    // we can easily share data between tasks or even jobs
    $`echo "files from another task: ${files}"`;
  });
});
```