const LOCAL_URLS: Record<string, string> = {
  ollama: "http://localhost:11434",
  vllm: "http://localhost:8000",
  lmstudio: "http://localhost:1234",
};

export interface DiscoveredModel {
  id: string;
  name: string;
}

export async function discoverModels(
  provider: string,
  baseUrl?: string
): Promise<DiscoveredModel[]> {
  const base = baseUrl || LOCAL_URLS[provider];
  if (!base) return [];

  try {
    if (provider === "ollama") {
      const res = await fetch(`${base}/api/tags`);
      if (!res.ok) return [];
      const data = await res.json();
      return (data.models ?? []).map((m: { name: string }) => ({
        id: m.name,
        name: m.name,
      }));
    } else {
      // OpenAI-compatible /v1/models
      const res = await fetch(`${base}/v1/models`);
      if (!res.ok) return [];
      const data = await res.json();
      return (data.data ?? []).map((m: { id: string }) => ({
        id: m.id,
        name: m.id,
      }));
    }
  } catch {
    return [];
  }
}
