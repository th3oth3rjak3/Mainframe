import { ThemeProvider } from "./components/providers/theme-provider";
import { Route, Routes } from "react-router-dom";
import Layout from "@/components/layout/layout";
import { Toaster } from "sonner";
import Dashboard from "@/pages/dashboard";

function App() {
  return (
    <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
      <Toaster richColors />
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
        </Routes>
      </Layout>
    </ThemeProvider>
  );
}

export default App;
