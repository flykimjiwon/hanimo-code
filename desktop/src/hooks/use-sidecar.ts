import { useEffect } from "react";
import {
  startSidecar,
  stopSidecar,
  sendPrompt,
  onSidecarEvent,
  type SidecarEvent,
} from "../lib/ipc";

export function useSidecar({
  onEvent,
}: {
  onEvent: (event: SidecarEvent) => void;
}) {
  useEffect(() => {
    let unlisten: (() => void) | undefined;

    const init = async () => {
      await startSidecar();
      unlisten = await onSidecarEvent(onEvent);
    };

    init().catch(console.error);

    return () => {
      unlisten?.();
      stopSidecar().catch(console.error);
    };
  }, []);

  const send = async (content: string) => {
    await sendPrompt(content);
  };

  return { send };
}
