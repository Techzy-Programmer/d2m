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
