import { ThemeProvider } from "./components/theme-provider";
import { Route, Routes } from "react-router-dom";
import Layout from "@/layout";
import Home from "@/components/home";

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </Layout>
    </ThemeProvider>
  );
}

export default App;
