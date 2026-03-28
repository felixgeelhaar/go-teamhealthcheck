import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { useWebSocket } from '../hooks/useWebSocket'
import { StatusBadge } from '../components/StatusBadge'
import { TrendChart } from '../components/TrendChart'
import { HealthAlerts } from '../components/HealthAlerts'
import { formatDate } from '../utils'
import type { Team, HealthCheck, WSEvent } from '../types'

interface TrendsData {
  team_id: string
  sessions: { id: string; name: string; date: string; avg_score: number; voters: number; status: string }[]
  trends: { MetricName: string; Sessions: { HealthCheckID: string; HealthCheckName: string; Score: number; Date: string }[]; Tendency: 'improving' | 'stable' | 'declining'; Delta: number }[]
}

export function TeamSelector() {
  const [selectedTeam, setSelectedTeam] = useState<string>('')
  const { data: teams, loading: teamsLoading } = useApi<Team[]>('/api/teams')
  const { data: healthchecks, loading: hcLoading, refetch } = useApi<HealthCheck[]>(
    selectedTeam ? `/api/healthchecks?team_id=${selectedTeam}` : '/api/healthchecks'
  )
  const { data: trends, refetch: refetchTrends } = useApi<TrendsData>(
    selectedTeam ? `/api/teams/${selectedTeam}/trends` : null
  )

  useWebSocket((event: WSEvent) => {
    if (['healthcheck_created', 'healthcheck_status_changed', 'healthcheck_deleted', 'vote_submitted'].includes(event.type)) {
      refetch()
      refetchTrends()
    }
  })

  const loading = teamsLoading || hcLoading

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '32px', flexWrap: 'wrap', gap: '16px' }}>
        <div>
          <h1 style={{ fontSize: '28px', fontWeight: 800, marginBottom: '4px' }}>Dashboard</h1>
          <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>
            Monitor team health and track improvements
          </p>
        </div>
        <div className="action-bar" style={{ margin: 0 }}>
          {selectedTeam && (
            <Link to={`/healthcheck/new/${selectedTeam}`} className="btn btn-primary">
              + New Health Check
            </Link>
          )}
          <Link to="/templates/new" className="btn btn-secondary">
            + New Template
          </Link>
        </div>
      </div>

      <div className="form-group">
        <label className="form-label">Team</label>
        <select
          className="form-select"
          value={selectedTeam}
          onChange={(e) => setSelectedTeam(e.target.value)}
          style={{ maxWidth: '300px' }}
        >
          <option value="">All teams</option>
          {teams?.map(t => (
            <option key={t.ID} value={t.ID}>{t.Name}</option>
          ))}
        </select>
      </div>

      {selectedTeam && <HealthAlerts teamId={selectedTeam} />}

      {trends && trends.sessions && trends.sessions.length > 0 && (
        <TrendChart sessions={trends.sessions} trends={trends.trends || []} />
      )}

      {loading && (
        <div className="loading">
          <div className="loading-spinner" />
          Loading...
        </div>
      )}

      {!loading && healthchecks && healthchecks.length === 0 && (
        <div className="empty-state">
          <div className="empty-state-icon">---</div>
          <div className="empty-state-text">
            No health checks yet.{' '}
            {selectedTeam
              ? 'Create one to get started.'
              : 'Select a team and create your first health check.'}
          </div>
        </div>
      )}

      {!loading && healthchecks && healthchecks.length > 0 && (
        <div className="hc-grid">
          {healthchecks.map(hc => {
            const team = teams?.find(t => t.ID === hc.TeamID)
            return (
              <Link
                key={hc.ID}
                to={`/healthcheck/${hc.ID}`}
                className="glass-card glass-card-interactive"
              >
                <div className="hc-card">
                  <div className="hc-card-header">
                    <div>
                      <div className="hc-card-name">{hc.Name}</div>
                      <div className="hc-card-date">{formatDate(hc.CreatedAt)}</div>
                    </div>
                    <StatusBadge status={hc.Status} />
                  </div>
                  <div className="hc-card-footer">
                    {team && <span>{team.Name}</span>}
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
