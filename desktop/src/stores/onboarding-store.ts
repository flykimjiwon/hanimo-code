import { create } from "zustand";
import { persist } from "zustand/middleware";

export const PROVIDERS = [
  { id: "openai", name: "OpenAI", models: ["gpt-4o", "gpt-4o-mini", "gpt-4.1", "gpt-4.1-mini"] },
  { id: "anthropic", name: "Anthropic", models: ["claude-sonnet-4-20250514", "claude-haiku-4-20250414", "claude-opus-4-20250514"] },
  { id: "google", name: "Google", models: ["gemini-2.5-flash", "gemini-2.5-pro"] },
  { id: "deepseek", name: "DeepSeek", models: ["deepseek-chat", "deepseek-coder", "deepseek-reasoner"] },
  { id: "groq", name: "Groq", models: ["qwen-qwq-32b", "llama-3.3-70b-versatile"] },
  { id: "ollama", name: "Ollama (Local)", models: [] },
];

interface OnboardingState {
  completed: boolean;
  provider: string;
  apiKey: string;
  model: string;
  step: number;
  setProvider: (provider: string) => void;
  setApiKey: (key: string) => void;
  setModel: (model: string) => void;
  nextStep: () => void;
  prevStep: () => void;
  complete: () => void;
  reset: () => void;
}

export const useOnboardingStore = create<OnboardingState>()(
  persist(
    (set) => ({
      completed: false,
      provider: "",
      apiKey: "",
      model: "",
      step: 0,
      setProvider: (provider) => set({ provider }),
      setApiKey: (key) => set({ apiKey: key }),
      setModel: (model) => set({ model }),
      nextStep: () => set((state) => ({ step: state.step + 1 })),
      prevStep: () => set((state) => ({ step: Math.max(0, state.step - 1) })),
      complete: () => set({ completed: true }),
      reset: () => set({ completed: false, provider: "", apiKey: "", model: "", step: 0 }),
    }),
    {
      name: "hanimo-onboarding",
    }
  )
);
