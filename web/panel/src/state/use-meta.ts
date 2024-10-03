import { create } from 'zustand'

type Meta = {
  pageTitle: string
  
  setPageTitle: (pageTitle: string) => void
}

export const useMeta = create<Meta>((set) => ({
  setPageTitle: (pageTitle) => set({ pageTitle }),

  pageTitle: 'Auth',
}));
