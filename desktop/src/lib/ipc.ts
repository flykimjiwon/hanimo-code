import { invoke } from "@tauri-apps/api/core";
import { listen, type UnlistenFn } from "@tauri-apps/api/event";

export interface SidecarEvent {
  type: "text" | "tool-call" | "tool-result" | "done" | "error";
  data: unknown;
  timestamp: string;
}

export interface SidecarConfig {
  provider: string;
  model: string;
  apiKey: string;
  role?: string;
  baseUrl?: string;
}

export async function startSidecar(config: SidecarConfig): Promise<string> {
  return invoke<string>("start_sidecar", {
    provider: config.provider,
    model: config.model,
    apiKey: config.apiKey,
    role: config.role,
    baseUrl: config.baseUrl,
  });
}

export async function sendPrompt(content: string): Promise<void> {
  return invoke<void>("send_prompt", { content });
}

export async function stopSidecar(): Promise<void> {
  return invoke<void>("stop_sidecar");
}

export async function onSidecarEvent(
  callback: (event: SidecarEvent) => void
): Promise<UnlistenFn> {
  return listen<string>("sidecar-event", (tauriEvent) => {
    try {
      const parsed = JSON.parse(tauriEvent.payload) as SidecarEvent;
      callback(parsed);
    } catch {
      // ignore non-JSON lines
    }
  });
}
