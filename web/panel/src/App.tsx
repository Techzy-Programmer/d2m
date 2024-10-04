import '@mantine/core/styles.css';
import '@mantine/dropzone/styles.css';
import '@mantine/notifications/styles.css';

import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { isFetchSuccess, showToast } from './utils/general';
import { Notifications } from '@mantine/notifications';
import { MantineProvider } from "@mantine/core"
import { useAuth } from './state/use-auth';
import { useMeta } from './state/use-meta';
import { MetaResp } from './utils/types';
import Header from "./components/Header"
import useFetch from './hooks/useFetch'
import Loading from './pages/Loading';
import { theme } from "./data/theme"
import { useEffect } from 'react';
import Home from './pages/Home';
import Auth from "./pages/Auth"
import "./App.css"

function App() {
  const { loggedIn, loading, setLoading, setLoggedIn } = useAuth();
  const { setPageTitle, setMetadata } = useMeta();
  const fetchData = useFetch();

  useEffect(() => {
    setLoading(true);
    const abt = new AbortController();

    (async () => {
      const { signal } = abt;
      const fetch = await fetchData<MetaResp>('/api/mg/meta', { signal });
      setLoading(false);

      if (!isFetchSuccess(fetch)) {
        return showToast({
          message: fetch.error || 'Something went wrong',
          title: 'Fetch Failed',
          status: 'issue'
        });
      }

      const { code, data } = fetch;

      if (code !== 200) {
        return showToast({
          title: 'Message From Server',
          message: data.message,
          status: 'warn'
        });
      }

      setMetadata(data.meta);
      setPageTitle("Home");
      setLoggedIn(true);

      showToast({
        title: 'Welcome back',
        message: 'You are now logged in',
        autoClose: 2000,
        status: 'ok'
      });
    })();

    return () => abt.abort();
  }, []);

  const homeNavigation = loading ? <Loading />
    : (loggedIn ? <Home /> : <Navigate to="/auth" />);

  return (
    <MantineProvider defaultColorScheme='auto' theme={theme}>
      <Header />
      <Notifications />

      <BrowserRouter>
        <Routes>
          <Route path='/'>
            <Route index element={homeNavigation} />
            
            <Route path="/auth" element={loading ? <Loading />
              : loggedIn ? <Navigate to="/" /> : <Auth />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </MantineProvider>
  )
}

export default App
