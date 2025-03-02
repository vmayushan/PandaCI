import { $, docker, job } from "@pandaci/workflow";

docker("ubuntu:latest", () => {
  $`echo "Hello, World!"`;
});
