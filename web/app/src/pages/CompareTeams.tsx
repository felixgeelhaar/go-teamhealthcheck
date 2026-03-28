import { Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { StatusBadge } from '../components/StatusBadge'
import { scoreColorClass, formatDate } from '../utils'

interface TeamComparison {
  team_id: string
  team_name: string
  healthcheck_name: string
  healthcheck_id: string
  avg_score: number
  voters: number
  date: string
  status: string
}

interface CompareResponse {
  teams: TeamComparison[]
}

export function CompareTeams() {
  const { data, loading } = useApi<CompareResponse>('/api/compare')

  const teams = (data?.teams || [])
    .slice()
    .sort((a, b) => b.avg_score - a.avg_score)

  return (
    <div>
      <Link to="/" className="back-link">
        {'\u2190'} Back to dashboard
      </Link>

      <div style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: '28px', fontWeight: 800, marginBottom: '4px', display: 'flex', alignItems: 'center', gap: '10px' }}>
          <span>{'\uD83C\uDFE2'}</span> Organization Health
        </h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>
          Compare health check results across teams
        </p>
      </div>

      {loading && (
        <div className="loading">
          <div className="loading-spinner" />
          Loading...
        </div>
      )}

      {!loading && teams.length === 0 && (
        <div className="empty-state">
          <div className="empty-state-icon">{'\uD83C\uDF10'}</div>
          <div className="empty-state-text">
            No health check data available
          </div>
        </div>
      )}

      {!loading && teams.length > 0 && (
        <div className="hc-grid">
          {teams.map(team => {
            const colorClass = scoreColorClass(team.avg_score)
            return (
              <Link
                key={team.healthcheck_id}
                to={`/healthcheck/${team.healthcheck_id}`}
                className="glass-card glass-card-interactive"
              >
                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <div>
                      <div style={{ fontWeight: 700, fontSize: '16px', color: 'var(--text-primary)' }}>
                        {team.team_name}
                      </div>
                      <div style={{ fontSize: '13px', color: 'var(--text-tertiary)', marginTop: '2px' }}>
                        {team.healthcheck_name}
                      </div>
                    </div>
                    <div
                      className={colorClass}
                      style={{
                        fontFamily: 'var(--font-mono)',
                        fontSize: '32px',
                        fontWeight: 900,
                        lineHeight: 1,
                      }}
                    >
                      {team.avg_score.toFixed(1)}
                    </div>
                  </div>

                  <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flexWrap: 'wrap' }}>
                    <StatusBadge status={team.status} />
                    <span style={{ fontSize: '13px', color: 'var(--text-tertiary)', fontFamily: 'var(--font-mono)' }}>
                      {team.voters} voter{team.voters !== 1 ? 's' : ''}
                    </span>
                    <span style={{ fontSize: '13px', color: 'var(--text-tertiary)' }}>
                      {formatDate(team.date)}
                    </span>
                  </div>
                </div>
              </Link>
            )
          })}
        </div>
      )}
    </div>
  )
}
