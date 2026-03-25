import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { useWebSocket } from '../hooks/useWebSocket'
import type { Team, HealthCheck, WSEvent } from '../types'

export function TeamSelector() {
  const [selectedTeam, setSelectedTeam] = useState<string>('')
  const { data: teams } = useApi<Team[]>('/api/teams')
  const { data: healthchecks, refetch } = useApi<HealthCheck[]>(
    selectedTeam ? `/api/healthchecks?team_id=${selectedTeam}` : '/api/healthchecks'
  )

  useWebSocket((event: WSEvent) => {
    if (['healthcheck_created', 'healthcheck_status_changed', 'healthcheck_deleted'].includes(event.type)) {
      refetch()
    }
  })

  const statusBadge = (status: string) => {
    const colors: Record<string, { bg: string; text: string }> = {
      open: { bg: '#dcfce7', text: '#166534' },
      closed: { bg: '#fef9c3', text: '#854d0e' },
      archived: { bg: '#f3f4f6', text: '#6b7280' },
    }
    const c = colors[status] || colors.archived
    return (
      <span style={{
        padding: '2px 8px',
        borderRadius: '12px',
        fontSize: '12px',
        fontWeight: 600,
        backgroundColor: c.bg,
        color: c.text,
      }}>
        {status}
      </span>
    )
  }

  return (
    <div>
      <h1 style={{ fontSize: '24px', fontWeight: 700, marginBottom: '24px' }}>
        Health Check Dashboard
      </h1>

      <div style={{ marginBottom: '24px' }}>
        <label style={{ fontSize: '14px', fontWeight: 600, display: 'block', marginBottom: '6px' }}>
          Team
        </label>
        <select
          value={selectedTeam}
          onChange={(e) => setSelectedTeam(e.target.value)}
          style={{
            padding: '8px 12px',
            borderRadius: '6px',
            border: '1px solid #d1d5db',
            fontSize: '14px',
            minWidth: '200px',
          }}
        >
          <option value="">All teams</option>
          {teams?.map(t => (
            <option key={t.ID} value={t.ID}>{t.Name}</option>
          ))}
        </select>
      </div>

      <div style={{
        display: 'grid',
        gap: '12px',
      }}>
        {healthchecks?.map(hc => (
          <Link
            key={hc.ID}
            to={`/healthcheck/${hc.ID}`}
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              padding: '16px',
              backgroundColor: 'white',
              borderRadius: '8px',
              border: '1px solid #e5e7eb',
              textDecoration: 'none',
              color: 'inherit',
              transition: 'box-shadow 0.15s',
            }}
            onMouseEnter={e => (e.currentTarget.style.boxShadow = '0 2px 8px rgba(0,0,0,0.08)')}
            onMouseLeave={e => (e.currentTarget.style.boxShadow = 'none')}
          >
            <div>
              <div style={{ fontWeight: 600, fontSize: '16px' }}>{hc.Name}</div>
              <div style={{ fontSize: '13px', color: '#6b7280', marginTop: '4px' }}>
                {new Date(hc.CreatedAt).toLocaleDateString()}
              </div>
            </div>
            {statusBadge(hc.Status)}
          </Link>
        ))}
        {healthchecks?.length === 0 && (
          <div style={{ color: '#9ca3af', padding: '24px', textAlign: 'center' }}>
            No health checks yet.
          </div>
        )}
      </div>
    </div>
  )
}
