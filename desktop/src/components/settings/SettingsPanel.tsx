import { useState } from "react";
import { useThemeStore } from "../../stores/theme-store";
import { useOnboardingStore, PROVIDERS } from "../../stores/onboarding-store";
import { discoverModels, DiscoveredModel } from "../../lib/model-discovery";

export function SettingsPanel() {
  const { theme } = useThemeStore();
  const c = theme.colors;
  const { provider, apiKey, model, baseUrl, setProvider, setApiKey, setModel, setBaseUrl, reset } =
    useOnboardingStore();

  const [localProvider, setLocalProvider] = useState(provider);
  const [localApiKey, setLocalApiKey] = useState(apiKey);
  const [localModel, setLocalModel] = useState(model);
  const [localBaseUrl, setLocalBaseUrl] = useState(baseUrl);

  const [discoveredModels, setDiscoveredModels] = useState<DiscoveredModel[]>([]);
  const [refreshing, setRefreshing] = useState(false);
  const [saved, setSaved] = useState(false);

  const selectedProvider = PROVIDERS.find((p) => p.id === localProvider);
  const needsApiKey = selectedProvider?.needsApiKey ?? true;
  const needsBaseUrl = (selectedProvider as { needsBaseUrl?: boolean } | undefined)?.needsBaseUrl ?? false;
  const availableModels = selectedProvider?.models ?? [];
  const isLocalProvider = selectedProvider && !selectedProvider.needsApiKey;

  function handleProviderChange(newProvider: string) {
    setLocalProvider(newProvider);
    setDiscoveredModels([]);
    const p = PROVIDERS.find((x) => x.id === newProvider);
    if (p && p.models.length > 0) {
      setLocalModel(p.models[0]);
    } else {
      setLocalModel("");
    }
  }

  async function handleRefreshModels() {
    setRefreshing(true);
    const models = await discoverModels(localProvider, localBaseUrl || undefined);
    setDiscoveredModels(models);
    if (models.length > 0 && !localModel) {
      setLocalModel(models[0].id);
    }
    setRefreshing(false);
  }

  function handleSave() {
    setProvider(localProvider);
    setApiKey(localApiKey);
    setModel(localModel);
    setBaseUrl(localBaseUrl);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  }

  const inputStyle = {
    background: c.bgSecondary,
    color: c.text,
    border: `1px solid ${c.border}`,
    borderRadius: 6,
    padding: "6px 10px",
    width: "100%",
    fontSize: 13,
    outline: "none",
  };

  const labelStyle = {
    fontSize: 11,
    fontWeight: 600,
    color: c.textMuted,
    textTransform: "uppercase" as const,
    letterSpacing: "0.05em",
    marginBottom: 4,
    display: "block",
  };

  const showDiscoveredSelect = isLocalProvider && discoveredModels.length > 0;
  const showStaticSelect = !isLocalProvider && availableModels.length > 0;

  return (
    <div
      className="flex flex-col flex-shrink-0 overflow-y-auto"
      style={{ width: 260, background: c.sidebarBg, borderRight: `1px solid ${c.border}` }}
    >
      <div style={{ padding: "12px 16px", borderBottom: `1px solid ${c.border}` }}>
        <span style={{ fontSize: 13, fontWeight: 600, color: c.text }}>Settings</span>
      </div>

      <div className="flex flex-col gap-4" style={{ padding: 16 }}>
        <div>
          <label style={labelStyle}>Provider</label>
          <select
            value={localProvider}
            onChange={(e) => handleProviderChange(e.target.value)}
            style={inputStyle}
          >
            {PROVIDERS.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>
        </div>

        {needsApiKey && (
          <div>
            <label style={labelStyle}>API Key</label>
            <input
              type="password"
              value={localApiKey}
              onChange={(e) => setLocalApiKey(e.target.value)}
              placeholder="sk-..."
              style={inputStyle}
            />
          </div>
        )}

        <div>
          <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: 4 }}>
            <label style={{ ...labelStyle, marginBottom: 0 }}>Model</label>
            {isLocalProvider && (
              <button
                onClick={handleRefreshModels}
                disabled={refreshing}
                style={{
                  background: "transparent",
                  color: c.textSecondary,
                  border: `1px solid ${c.border}`,
                  borderRadius: 4,
                  padding: "2px 8px",
                  fontSize: 11,
                  cursor: refreshing ? "not-allowed" : "pointer",
                  opacity: refreshing ? 0.6 : 1,
                }}
              >
                {refreshing ? "Loading..." : "Refresh Models"}
              </button>
            )}
          </div>
          {showDiscoveredSelect ? (
            <select
              value={localModel}
              onChange={(e) => setLocalModel(e.target.value)}
              style={inputStyle}
            >
              {discoveredModels.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.name}
                </option>
              ))}
            </select>
          ) : showStaticSelect ? (
            <select
              value={localModel}
              onChange={(e) => setLocalModel(e.target.value)}
              style={inputStyle}
            >
              {availableModels.map((m) => (
                <option key={m} value={m}>
                  {m}
                </option>
              ))}
            </select>
          ) : (
            <input
              type="text"
              value={localModel}
              onChange={(e) => setLocalModel(e.target.value)}
              placeholder="model name"
              style={inputStyle}
            />
          )}
        </div>

        {needsBaseUrl && (
          <div>
            <label style={labelStyle}>Base URL</label>
            <input
              type="text"
              value={localBaseUrl}
              onChange={(e) => setLocalBaseUrl(e.target.value)}
              placeholder="https://..."
              style={inputStyle}
            />
          </div>
        )}

        <button
          onClick={handleSave}
          style={{
            background: c.accent,
            color: "#fff",
            border: "none",
            borderRadius: 6,
            padding: "8px 12px",
            fontSize: 13,
            fontWeight: 600,
            cursor: "pointer",
            width: "100%",
          }}
        >
          {saved ? "Saved! Reconnecting..." : "Save & Reconnect"}
        </button>

        <div style={{ borderTop: `1px solid ${c.border}`, paddingTop: 16 }}>
          <button
            onClick={reset}
            style={{
              background: "transparent",
              color: c.textMuted,
              border: `1px solid ${c.border}`,
              borderRadius: 6,
              padding: "8px 12px",
              fontSize: 13,
              cursor: "pointer",
              width: "100%",
            }}
          >
            Reset Onboarding
          </button>
        </div>
      </div>
    </div>
  );
}
