import { useEffect, useState } from 'react'
import { ShieldCheck, ShieldAlert, ShieldQuestion, Shield, Sparkles } from 'lucide-react'

// Phase 14b — Permissions panel skeleton.
//
// 5-mode toggle UI only. Real enforcement (yaml rules, credential scrub,
// sandbox) lives in a future phase — for now this lets the user pick a
// mode that other panels can read via window.localStorage. Cycling with
// Shift+Tab matches the design v1 mock.

type Mode = 'allow' | 'ask' | 'deny' | 'dangerous-block' | 'learn'

const MODES: { id: Mode; label: string; sub: string; Icon: any; tone: 'success' | 'accent' | 'error' | 'muted' }[] = [
  { id: 'allow',           label: 'Allow All',        sub: '모든 도구 자동 승인 (개발 환경 전용)',         Icon: ShieldCheck,   tone: 'muted' },
  { id: 'ask',             label: 'Ask Each',         sub: '도구 호출 직전 매번 확인 (기본값)',            Icon: ShieldQuestion, tone: 'accent' },
  { id: 'deny',            label: 'Read-Only',        sub: '쓰기 도구 차단, 읽기·검색만 허용',             Icon: Shield,        tone: 'muted' },
  { id: 'dangerous-block', label: 'Block Dangerous',  sub: 'rm -rf · sudo · git push --force 등 자동 차단', Icon: ShieldAlert,   tone: 'error' },
  { id: 'learn',           label: 'Learning',         sub: '사용자 응답을 yaml 룰로 누적 (장기)',         Icon: Sparkles,      tone: 'success' },
]

const STORAGE_KEY = 'hanimo-permissions-mode'

export default function PermissionsPanel() {
  const [mode, setMode] = useState<Mode>(() => {
    const stored = (typeof window !== 'undefined' ? window.localStorage.getItem(STORAGE_KEY) : null) as Mode | null
    return stored && MODES.some(m => m.id === stored) ? stored : 'ask'
  })

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEY, mode)
    }
  }, [mode])

  // Shift+Tab cycles modes when this panel is mounted (matches mock UX).
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.shiftKey && e.key === 'Tab' && (e.target as HTMLElement)?.tagName !== 'INPUT' && (e.target as HTMLElement)?.tagName !== 'TEXTAREA') {
        e.preventDefault()
        setMode(prev => {
          const idx = MODES.findIndex(m => m.id === prev)
          return MODES[(idx + 1) % MODES.length].id
        })
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [])

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
        <ShieldCheck size={14} style={{ color: 'var(--accent)' }} />
        <span style={{ fontSize: 11, fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--fg-muted)' }}>
          Permissions
        </span>
        <span style={{
          marginLeft: 'auto',
          fontSize: 9, fontWeight: 700,
          background: 'var(--accent-glow)', color: 'var(--accent)',
          padding: '1px 6px', borderRadius: 3, textTransform: 'uppercase', letterSpacing: '0.04em',
        }}>
          UI 골격
        </span>
      </div>

      <div style={{ flex: 1, overflowY: 'auto', padding: '10px 0 14px' }}>
        <div style={{ padding: '0 14px 8px', fontSize: 11, color: 'var(--fg-dim)', lineHeight: 1.55 }}>
          현재 모드를 선택하세요. <kbd style={{
            fontFamily: 'var(--font-code)', fontSize: 10,
            background: 'var(--bg-active)', padding: '1px 5px',
            borderRadius: 3, color: 'var(--fg-muted)',
          }}>Shift+Tab</kbd>로 순환할 수 있습니다.
        </div>

        {MODES.map(m => {
          const active = m.id === mode
          const toneColor =
            m.tone === 'success' ? 'var(--success, #3fb950)' :
            m.tone === 'error'   ? 'var(--error, #f85149)' :
            m.tone === 'accent'  ? 'var(--accent)' :
                                   'var(--fg-muted)'
          return (
            <button
              key={m.id}
              onClick={() => setMode(m.id)}
              style={{
                width: '100%',
                textAlign: 'left',
                display: 'flex',
                alignItems: 'flex-start',
                gap: 10,
                padding: '8px 14px',
                background: active ? 'var(--bg-hover)' : 'none',
                border: 'none',
                borderLeft: active ? '2px solid var(--accent)' : '2px solid transparent',
                cursor: 'pointer',
                fontFamily: 'var(--font-ui)',
                transition: 'background 0.1s, border-color 0.1s',
              }}
              onMouseEnter={e => {
                if (!active) (e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-hover)'
              }}
              onMouseLeave={e => {
                if (!active) (e.currentTarget as HTMLButtonElement).style.background = 'none'
              }}
            >
              <m.Icon size={14} style={{ color: toneColor, marginTop: 2, flexShrink: 0 }} />
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                  <span style={{
                    color: active ? 'var(--fg-primary)' : 'var(--fg-secondary)',
                    fontSize: 12.5,
                    fontWeight: active ? 600 : 500,
                  }}>
                    {m.label}
                  </span>
                  {active && (
                    <span style={{
                      fontSize: 9, fontWeight: 700,
                      background: toneColor, color: '#000',
                      padding: '1px 5px', borderRadius: 3, textTransform: 'uppercase', letterSpacing: '0.04em',
                    }}>
                      Active
                    </span>
                  )}
                </div>
                <div style={{ fontSize: 11, color: 'var(--fg-muted)', lineHeight: 1.45, marginTop: 2 }}>
                  {m.sub}
                </div>
              </div>
            </button>
          )
        })}

        <div style={{
          margin: '12px 12px 0',
          padding: '8px 10px',
          background: 'var(--bg-hover)',
          borderRadius: 4,
          borderLeft: '2px solid var(--accent)',
          fontSize: 11,
          color: 'var(--fg-muted)',
          lineHeight: 1.5,
        }}>
          모드 선택은 localStorage에만 저장됩니다. 실제 도구 호출 차단/검사 로직은
          다음 Phase에서 추가될 예정 — 현재는 UI 흐름 + 사용자 의도 기록만.
        </div>
      </div>
    </div>
  )
}
