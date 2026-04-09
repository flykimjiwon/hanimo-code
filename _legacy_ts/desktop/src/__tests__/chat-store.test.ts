import { describe, it, expect, beforeEach } from "vitest";
import { useChatStore } from "../stores/chat-store";

describe("ChatStore", () => {
  beforeEach(() => {
    useChatStore.setState({
      messages: [],
      isStreaming: false,
      streamingContent: "",
    });
  });

  it("adds a user message", () => {
    useChatStore.getState().addMessage({ role: "user", content: "hello" });
    const msgs = useChatStore.getState().messages;
    expect(msgs).toHaveLength(1);
    expect(msgs[0].role).toBe("user");
    expect(msgs[0].content).toBe("hello");
  });

  it("manages streaming state", () => {
    const store = useChatStore.getState();
    store.setStreaming(true);
    store.appendStreamingContent("hello ");
    store.appendStreamingContent("world");
    expect(useChatStore.getState().isStreaming).toBe(true);
    expect(useChatStore.getState().streamingContent).toBe("hello world");
  });

  it("clears messages", () => {
    useChatStore.getState().addMessage({ role: "user", content: "test" });
    useChatStore.getState().clear();
    expect(useChatStore.getState().messages).toHaveLength(0);
  });
});
