import { create } from 'zustand'

type Auth = {
  loading: boolean
  loggedIn: boolean
  
  setLoading: (loading: boolean) => void
  setLoggedIn: (loggedIn: boolean) => void
}

export const useAuth = create<Auth>((set) => ({
  setLoggedIn: (loggedIn) => set({ loggedIn }),
  setLoading: (loading) => set({ loading }),

  loggedIn: false,
  loading: true,
}));
