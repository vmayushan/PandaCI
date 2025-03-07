import { Conclusion as ProtoConclusion } from "@pandaci/proto"

export type Conclusion = "success" | "failure" | "skipped";

export function protoConclusionToConclusion(proto: ProtoConclusion): Conclusion {
  switch (proto) {
    case ProtoConclusion.SUCCESS:
      return "success";
    case ProtoConclusion.FAILURE:
      return "failure";
    case ProtoConclusion.SKIPPED: 
      return "skipped";
    default:
      throw new Error(`Unknown conclusion: ${proto}`);
  }
}