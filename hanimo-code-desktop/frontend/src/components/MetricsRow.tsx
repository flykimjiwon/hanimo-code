interface Metric {
  label: string
  value: string
  sub?: string
  subColor?: 'up' | 'accent' | 'muted'
}

interface Props {
  contextPct?: number
  contextTokens?: number
  contextMax?: number
  cacheHitPct?: number
  cacheSavedUsd?: number
  iter?: number
  iterMax?: number
  iterLabel?: string
  provider?: string
  tier?: 'T1' | 'T2' | 'T3'
}

/**
 * MetricsRow — 4열 지표 스트립 for the right panel.
 * Phase 2 scaffold: accepts optional props, shows mock values when missing so
 * the component renders immediately. Real metrics pipe (Go GetMetrics bind)
 * is Phase 3+.
 */
export default function MetricsRow(props: Props) {
  const metrics: Metric[] = [
    {
      label: 'Context',
      value: props.contextPct !== undefined ? `${props.contextPct}%` : '·',
      sub:
        props.contextTokens !== undefined && props.contextMax !== undefined
          ? `${props.contextTokens.toLocaleString()} / ${(props.contextMax / 1000).toFixed(0)}K`
          : '— / —',
      subColor: 'up',
    },
    {
      label: 'Cache hit',
      value: props.cacheHitPct !== undefined ? `${props.cacheHitPct}%` : '·',
      sub:
        props.cacheSavedUsd !== undefined
          ? `saved ~$${props.cacheSavedUsd.toFixed(2)}`
          : '—',
      subColor: 'accent',
    },
    {
      label: 'Iter',
      value:
        props.iter !== undefined && props.iterMax !== undefined
          ? `${props.iter} / ${props.iterMax}`
          : '—',
      sub: props.iterLabel ?? 'idle',
      subColor: 'up',
    },
    {
      label: 'Provider',
      value: props.provider ?? '—',
      sub: props.tier ? `${props.tier} certified` : 'unrated',
      subColor: 'accent',
    },
  ]

  return (
    <div
      style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(4, 1fr)',
        gap: 12,
        padding: '8px 14px',
        borderBottom: '1px solid var(--border)',
        fontSize: 10,
        color: 'var(--fg-muted)',
        fontFamily: 'var(--font-ui)',
      }}
    >
      {metrics.map(m => (
        <div key={m.label} style={{ display: 'flex', flexDirection: 'column', gap: 1, minWidth: 0 }}>
          <span>{m.label}</span>
          <b
            style={{
              color: m.label === 'Provider' ? 'var(--accent)' : 'var(--fg-primary)',
              fontSize: 12,
              fontWeight: 600,
              fontFamily: 'var(--font-ui)',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
          >
            {m.value}
          </b>
          <span
            style={{
              fontSize: 9,
              color:
                m.subColor === 'accent'
                  ? 'var(--accent)'
                  : m.subColor === 'up'
                  ? 'var(--success)'
                  : 'var(--fg-muted)',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
          >
            {m.sub}
          </span>
        </div>
      ))}
    </div>
  )
}
