import { useEffect } from "react";
import {
  startSidecar,
  stopSidecar,
  sendPrompt,
  onSidecarEvent,
  type SidecarEvent,
} from "../lib/ipc";
import { useOnboardingStore } from "../stores/onboarding-store";

export function useSidecar({
  onEvent,
}: {
  onEvent: (event: SidecarEvent) => void;
}) {
  const { provider, apiKey, model } = useOnboardingStore();

  useEffect(() => {
    let unlisten: (() => void) | undefined;

    const init = async () => {
      await startSidecar({ provider, model, apiKey, role: "hanimo" });
      unlisten = await onSidecarEvent(onEvent);
    };

    init().catch(console.error);

    return () => {
      unlisten?.();
      stopSidecar().catch(console.error);
    };
  }, [onEvent, provider, model, apiKey]);

  const send = async (content: string) => {
    await sendPrompt(content);
  };

  return { send };
}
