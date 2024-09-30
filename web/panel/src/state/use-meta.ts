import { create } from 'zustand'

type Meta = {
  loggedIn: boolean
  pageTitle: string
  
  setLoggedIn: (loggedIn: boolean) => void
  setPageTitle: (pageTitle: string) => void
}

export const useMeta = create<Meta>((set) => ({
  setPageTitle: (pageTitle) => set({ pageTitle }),
  setLoggedIn: (loggedIn) => set({ loggedIn }),

  pageTitle: 'Auth',
  loggedIn: false,
}));
