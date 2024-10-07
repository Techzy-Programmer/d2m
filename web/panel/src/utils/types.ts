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
