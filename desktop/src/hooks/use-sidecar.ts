import { useEffect, useState, useCallback } from "react";
import {
  startSidecar,
  stopSidecar,
  sendPrompt,
  onSidecarEvent,
  onSidecarStderr,
  type SidecarEvent,
} from "../lib/ipc";
import { useOnboardingStore } from "../stores/onboarding-store";
import { useChatStore } from "../stores/chat-store";
import { useEditorStore } from "../stores/editor-store";

export function useSidecar({
  onEvent,
  role = "hanimo",
}: {
  onEvent: (event: SidecarEvent) => void;
  role?: string;
}) {
  const { provider, apiKey, model, baseUrl } = useOnboardingStore();
  const { setConnectionStatus, setConnectionError } = useChatStore();
  const { projectRoot } = useEditorStore();
  const [retryCount, setRetryCount] = useState(0);

  useEffect(() => {
    let unlisten: (() => void) | undefined;
    let unlistenStderr: (() => void) | undefined;

    const init = async () => {
      try {
        setConnectionStatus("connecting");
        await startSidecar({
          provider,
          model,
          apiKey,
          role,
          baseUrl,
          projectDir: projectRoot ?? undefined,
        });
        unlisten = await onSidecarEvent(onEvent);
        unlistenStderr = await onSidecarStderr((msg) => {
          console.error("[sidecar stderr]", msg);
        });
        setConnectionStatus("connected");
      } catch (err) {
        setConnectionStatus("error");
        setConnectionError(err instanceof Error ? err.message : String(err));
      }
    };

    init();

    return () => {
      unlisten?.();
      unlistenStderr?.();
      stopSidecar().catch(console.error);
    };
  }, [onEvent, provider, model, apiKey, role, baseUrl, projectRoot, retryCount]);

  const send = async (content: string) => {
    await sendPrompt(content);
  };

  const retry = useCallback(() => {
    setRetryCount((n) => n + 1);
  }, []);

  return { send, retry };
}
