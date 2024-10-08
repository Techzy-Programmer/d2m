import { Clock, GitBranch, GitCommitVertical, Info, MessageCircle } from "lucide-react";
import { DeploymentDetailResp } from "../utils/types";

type DeploymentMeta = {
  key: keyof DeploymentDetailResp["meta"];
  icon: JSX.Element;
  hovText: string;
}

export const deploymentsMeta: DeploymentMeta[] = [
  {
    key: "StartAt",
    icon: <Clock />,
    hovText: "Starting point of deployment",
  },
  {
    key: "EndAt",
    icon: <Clock />,
    hovText: "Ending point of deployment",
  },
  {
    key: "CommitHash",
    icon: <GitCommitVertical />,
    hovText: "Commit hash of deployment",
  },
  {
    key: "CommitMsg",
    icon: <MessageCircle />,
    hovText: "Commit message of deployment",
  },
  {
    key: "Branch",
    icon: <GitBranch />,
    hovText: "Deployment branch from repository",
  },
  {
    key: "Status",
    icon: <Info />,
    hovText: "Deployment status",
  }
]
