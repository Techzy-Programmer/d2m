import { create } from 'zustand'
import { MetaResp } from '../utils/types'

type Meta = {
  pageTitle: string
  metadata: MetaResp["meta"]
  
  setMetadata: (metadata: MetaResp["meta"]) => void
  setPageTitle: (pageTitle: string) => void
}

export const useMeta = create<Meta>((set) => ({
  setPageTitle: (pageTitle) => set({ pageTitle }),
  setMetadata: (metadata) => set({ metadata }),

  metadata: {} as MetaResp["meta"],
  pageTitle: 'Auth',
}));
