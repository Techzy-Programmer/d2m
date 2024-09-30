
import '@mantine/core/styles.css';
import '@mantine/dropzone/styles.css';

import { MantineProvider } from "@mantine/core"
import Header from "./components/Header"
import { theme } from "./data/theme"
import Auth from "./pages/Auth"
import "./App.css"


function App() {
  return (
    <MantineProvider defaultColorScheme='auto' theme={theme}>
      <Header />
      <Auth />
    </MantineProvider>
  )
}

export default App
