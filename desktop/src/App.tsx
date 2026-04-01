import { Sidebar } from "./components/layout/Sidebar";
import { FileTree } from "./components/layout/FileTree";
import { EditorArea } from "./components/layout/EditorArea";
import { ChatPanel } from "./components/layout/ChatPanel";
import { useThemeStore } from "./stores/theme-store";

export function App() {
  const { theme } = useThemeStore();
  const c = theme.colors;

  return (
    <div className="h-screen flex overflow-hidden" style={{ background: c.bg, color: c.text }}>
      <Sidebar />
      <FileTree />
      <EditorArea />
      <ChatPanel />
    </div>
  );
}
