import { Link } from 'react-router-dom'
import { useParticipant } from '../hooks/useParticipant'

function getAvatarColor(name: string): string {
  const colors = ['#3b82f6', '#a855f7', '#ec4899', '#f97316', '#14b8a6', '#6366f1']
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

function getInitial(name: string): string {
  return name.charAt(0).toUpperCase()
}

interface Props {
  onChangeNameClick: () => void
}

export function Navbar({ onChangeNameClick }: Props) {
  const { name } = useParticipant()

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-brand">
        <span className="pulse-dot" />
        Health Check
      </Link>

      <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
        <Link to="/compare" className="btn btn-ghost btn-sm">
          {'\uD83C\uDFE2'} Compare Teams
        </Link>
        {name && (
          <div className="navbar-user" onClick={onChangeNameClick}>
            <div
              className="avatar"
              style={{ backgroundColor: getAvatarColor(name) }}
            >
              {getInitial(name)}
            </div>
            <span className="navbar-user-name">{name}</span>
          </div>
        )}
      </div>
    </nav>
  )
}
