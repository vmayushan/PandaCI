import { WorkflowService } from "@pandaci/proto";
import {
  type Client,
  createClient,
  createConnectTransport,
  type Interceptor,
} from "./deps.ts";
import { extendedEnv } from "./env.ts";

let workflowClient: Client<typeof WorkflowService> | undefined;

const auth: Interceptor = (next) => (req) => {
  req.header.set("Authorization", `Bearer ${extendedEnv.PANDACI_WORKFLOW_JWT}`);
  return  next(req);
};

export function getWorkflowClient(): Client<typeof WorkflowService> {
  if (workflowClient) {
    return workflowClient;
  }

  const transport = createConnectTransport({
    baseUrl: extendedEnv.PANDACI_WORKFLOW_GRPC_ADDRESS + "/grpc",
    interceptors: [auth],
  });

  workflowClient = createClient(WorkflowService, transport);

  return workflowClient;
}
