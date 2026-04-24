import { useState, useEffect } from 'react'
import {
  Files,
  Search,
  GitBranch,
  Play,
  TriangleAlert,
  BookOpen,
  Sparkle,
  PlugZap,
  Share2,
  History,
  Monitor,
  ShieldCheck,
  User,
  Settings,
} from 'lucide-react'
import { GetGitInfo } from '../../wailsjs/go/main/App'

interface Props {
  active: string
  onSelect: (panel: string) => void
}

type Item = {
  id: string
  icon: typeof Files
  label: string
  badge?: number
  dot?: boolean
}

export default function ActivityBar({ active, onSelect }: Props) {
  const [gitCount, setGitCount] = useState(0)

  useEffect(() => {
    GetGitInfo().then(info => setGitCount(info.changes?.length || 0)).catch(() => {})
    const interval = setInterval(() => {
      GetGitInfo().then(info => setGitCount(info.changes?.length || 0)).catch(() => {})
    }, 5000)
    return () => clearInterval(interval)
  }, [])

  const top: Item[] = [
    { id: 'files', icon: Files, label: 'Explorer' },
    { id: 'search', icon: Search, label: 'Search' },
    { id: 'git', icon: GitBranch, label: 'Git', badge: gitCount, dot: gitCount > 0 },
    { id: 'run', icon: Play, label: 'Run / Debug' },
    { id: 'problems', icon: TriangleAlert, label: 'LSP Problems' },
    { id: 'knowledge', icon: BookOpen, label: 'Knowledge Packs' },
    { id: 'skills', icon: Sparkle, label: 'Skills' },
    { id: 'mcp', icon: PlugZap, label: 'MCP Servers' },
    { id: 'subagents', icon: Share2, label: 'Subagents' },
    { id: 'sessions', icon: History, label: 'Sessions' },
    { id: 'webpreview', icon: Monitor, label: 'Web Preview' },
  ]

  const bottom: Item[] = [
    { id: 'permissions', icon: ShieldCheck, label: 'Permissions' },
    { id: 'account', icon: User, label: 'Account' },
    { id: 'settings', icon: Settings, label: 'Settings' },
  ]

  const renderIcon = (item: Item, isBottom = false) => (
    <button
      key={item.id}
      title={item.label}
      aria-label={item.label}
      onClick={() => onSelect(item.id)}
      style={{
        background: 'none',
        border: 'none',
        cursor: 'pointer',
        padding: 8,
        borderRadius: 8,
        color: active === item.id && !isBottom ? 'var(--accent)' : 'var(--fg-muted)',
        position: 'relative',
        transition: 'all 0.15s',
      }}
      onMouseEnter={e => {
        (e.currentTarget as HTMLButtonElement).style.color = 'var(--fg-primary)'
        ;(e.currentTarget as HTMLButtonElement).style.background = 'var(--bg-hover)'
      }}
      onMouseLeave={e => {
        (e.currentTarget as HTMLButtonElement).style.color =
          active === item.id && !isBottom ? 'var(--accent)' : 'var(--fg-muted)'
        ;(e.currentTarget as HTMLButtonElement).style.background = 'none'
      }}
    >
      {active === item.id && !isBottom && (
        <span
          style={{
            position: 'absolute',
            left: 0,
            top: 8,
            bottom: 8,
            width: 2,
            background: 'var(--accent)',
            borderRadius: '0 2px 2px 0',
          }}
        />
      )}
      <item.icon size={18} />
      {item.badge && item.badge > 0 ? (
        <span
          style={{
            position: 'absolute',
            top: 2,
            right: 2,
            background: 'var(--accent)',
            color: 'var(--accent-text, #fff)',
            fontSize: 8,
            fontWeight: 700,
            minWidth: 14,
            height: 14,
            borderRadius: 7,
            padding: '0 3px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontFamily: 'var(--font-ui)',
          }}
        >
          {item.badge > 9 ? '9+' : item.badge}
        </span>
      ) : item.dot ? (
        <span
          style={{
            position: 'absolute',
            top: 6,
            right: 6,
            width: 6,
            height: 6,
            borderRadius: 3,
            background: 'var(--success)',
          }}
        />
      ) : null}
    </button>
  )

  return (
    <aside
      style={{
        width: 48,
        background: 'var(--bg-activity)',
        borderRight: '1px solid var(--border)',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        padding: '8px 0',
        gap: 1,
        overflowY: 'auto',
      }}
    >
      {top.map(item => renderIcon(item))}
      <div style={{ flex: 1, minHeight: 8 }} />
      {bottom.map(item => renderIcon(item, true))}
    </aside>
  )
}
