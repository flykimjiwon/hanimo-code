import { useEffect, useState } from 'react'
import { Play, RefreshCw, Package, Hammer, Cog } from 'lucide-react'

interface RunTarget {
  name: string
  description: string
  command: string
  source: 'npm' | 'make' | 'go'
}

const sourceIcon = {
  npm:  Package,
  make: Hammer,
  go:   Cog,
} as const

export default function RunPanel() {
  const [targets, setTargets] = useState<RunTarget[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState('')

  const load = () => {
    setLoading(true)
    import('../../wailsjs/go/main/App').then(mod => {
      const fn = (mod as any).GetRunTargets
      if (typeof fn !== 'function') { setLoading(false); return }
      fn().then((res: RunTarget[] | null) => {
        setTargets(Array.isArray(res) ? res : [])
        setLoading(false)
      }).catch(() => setLoading(false))
    }).catch(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  function run(t: RunTarget) {
    import('../../wailsjs/go/main/App').then((mod: any) => {
      const fn = mod.RunTarget
      if (typeof fn !== 'function') {
        import('./Toast').then(toast => toast.showToast('RunTarget binding pending — wails dev', 'info'))
        return
      }
      fn(t.command).then(() => {
        import('./Toast').then(toast => toast.showToast(`▶ ${t.command}`, 'success'))
      }).catch((e: any) => {
        import('./Toast').then(toast => toast.showToast(`Run failed: ${e?.message || 'unknown'}`, 'info'))
      })
    }).catch(() => {})
  }

  const filtered = filter
    ? targets.filter(t => (t.name + ' ' + t.command + ' ' + t.description).toLowerCase().includes(filter.toLowerCase()))
    : targets

  // Group by source so the user sees npm/make/go as separate sections
  const groups: Record<string, RunTarget[]> = {}
  for (const t of filtered) (groups[t.source] ??= []).push(t)
  const orderedSources = ['npm', 'make', 'go'].filter(s => groups[s]?.length)

  return (
    <div style={{
      height: '100%',
      background: 'var(--bg-sidebar)',
      borderRight: '1px solid var(--border)',
      display: 'flex',
      flexDirection: 'column',
      overflow: 'hidden',
    }}>
      <div style={{
        padding: '10px 14px',
        borderBottom: '1px solid var(--border)',
        display: 'flex',
        alignItems: 'center',
        gap: 6,
      }}>
        <Play size={14} style={{ color: 'var(--accent)' }} />
        <span style={{ fontSize: 11, fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--fg-muted)' }}>
          Run / Debug
        </span>
        <span style={{ marginLeft: 'auto', fontSize: 10, color: 'var(--fg-dim)' }}>
          {targets.length}
        </span>
        <button onClick={load} title="Rescan project" aria-label="Refresh" style={{
          background: 'none', border: 'none', cursor: 'pointer', color: 'var(--fg-muted)',
          padding: 0, display: 'flex',
        }}>
          <RefreshCw size={12} />
        </button>
      </div>

      <input
        value={filter}
        onChange={e => setFilter(e.target.value)}
        placeholder="Filter targets…"
        style={{
          margin: '8px 12px',
          background: 'var(--bg-input)',
          border: '1px solid var(--border)',
          borderRadius: 4,
          padding: '4px 8px',
          color: 'var(--fg-primary)',
          fontSize: 12,
          fontFamily: 'var(--font-ui)',
          outline: 'none',
        }}
      />

      <div style={{ flex: 1, overflowY: 'auto', padding: '4px 0 14px' }}>
        {loading && (
          <div style={{ padding: '8px 14px', fontSize: 11, color: 'var(--fg-dim)' }}>Loading…</div>
        )}

        {!loading && targets.length === 0 && (
          <div style={{ padding: '12px 14px', fontSize: 11, color: 'var(--fg-dim)', lineHeight: 1.55 }}>
            실행 가능한 스크립트가 없습니다.<br />
            <code style={{
              fontSize: 10.5, background: 'var(--bg-active)', padding: '1px 4px', borderRadius: 3,
              color: 'var(--accent)',
            }}>package.json</code>의 <code style={{
              fontSize: 10.5, background: 'var(--bg-active)', padding: '1px 4px', borderRadius: 3,
              color: 'var(--accent)',
            }}>scripts</code>나 <code style={{
              fontSize: 10.5, background: 'var(--bg-active)', padding: '1px 4px', borderRadius: 3,
              color: 'var(--accent)',
            }}>Makefile</code>을 추가하세요.
          </div>
        )}

        {orderedSources.map(src => {
          const Icon = sourceIcon[src as keyof typeof sourceIcon]
          return (
            <div key={src}>
              <div style={{
                padding: '6px 14px 2px',
                fontSize: 9,
                fontWeight: 700,
                textTransform: 'uppercase',
                letterSpacing: '0.06em',
                color: 'var(--fg-dim)',
                display: 'flex',
                alignItems: 'center',
                gap: 5,
              }}>
                <Icon size={10} />
                {src}
                <span style={{ color: 'var(--fg-dim)', fontWeight: 500 }}>
                  · {groups[src].length}
                </span>
              </div>
              {groups[src].map(t => (
                <button
                  key={t.source + ':' + t.name}
                  onClick={() => run(t)}
                  title={t.command}
                  style={{
                    width: '100%',
                    textAlign: 'left',
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 2,
                    padding: '5px 14px',
                    background: 'none',
                    border: 'none',
                    borderLeft: '2px solid transparent',
                    cursor: 'pointer',
                    fontFamily: 'var(--font-ui)',
                    transition: 'background 0.1s, border-color 0.1s',
                  }}
                  onMouseEnter={e => {
                    (e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-hover)'
                    ;(e.currentTarget as HTMLButtonElement).style.borderLeftColor = 'var(--accent)'
                  }}
                  onMouseLeave={e => {
                    (e.currentTarget as HTMLButtonElement).style.background = 'none'
                    ;(e.currentTarget as HTMLButtonElement).style.borderLeftColor = 'transparent'
                  }}
                >
                  <span style={{
                    color: 'var(--fg-primary)', fontSize: 12.5, fontWeight: 500,
                    fontFamily: 'var(--font-code)',
                    overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                  }}>
                    {t.name}
                  </span>
                  <span style={{
                    fontSize: 10.5, color: 'var(--fg-muted)', lineHeight: 1.45,
                    overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
                  }}>
                    {t.description}
                  </span>
                </button>
              ))}
            </div>
          )
        })}
      </div>
    </div>
  )
}
