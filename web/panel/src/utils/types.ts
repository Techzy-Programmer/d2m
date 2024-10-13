type Resp = {
  message: string;
  ok: boolean;
}

export type MetaResp = Resp & {
  meta: {
    uptime: string
    webPort: string
    tcpPort: string
  }
}

type Deployment = {
  ID: number
  Repo: string
  Branch: string
  CommitHash: string
  CommitMsg: string
  StartAt: string
  EndAt: string
  Status: string
}

export type DeploymentsResp = Resp & {
  deployments: Deployment[]
}

type Log = {
  ID: number
  Level: number
  DeployID: number
  Timestamp: number
  Message: string
  Steps: string
  Title: string
}

export type DeploymentDetailResp = Resp & {
  meta: Deployment
  logs: Log[]
}
