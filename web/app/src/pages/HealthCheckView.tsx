import { useParams, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { useWebSocket } from '../hooks/useWebSocket'
import { MetricRow } from '../components/MetricRow'
import { VoteProgress } from '../components/VoteProgress'
import { StatusBadge } from '../components/StatusBadge'
import { RadarChart } from '../components/RadarChart'
import { DiscussionGuide } from '../components/DiscussionGuide'
import { ActionItems } from '../components/ActionItems'
import { scoreColorClass, getAvatarColor, getInitial } from '../utils'
import type { HealthCheckResults, PluginEntry, WSEvent } from '../types'

export function HealthCheckView() {
  const { id } = useParams<{ id: string }>()
  const { data, loading, refetch } = useApi<HealthCheckResults>(
    id ? `/api/healthchecks/${id}/results` : null
  )
  const { data: plugins } = useApi<PluginEntry[]>('/api/plugins')

  useWebSocket((event: WSEvent) => {
    if (event.healthcheck_id === id) {
      refetch()
    }
  })

  if (loading || !data) {
    return (
      <div className="loading">
        <div className="loading-spinner" />
        Loading...
      </div>
    )
  }

  const { healthcheck: hc, results, average_score, participants, participant_names, total_votes } = data
  const colorClass = scoreColorClass(average_score)

  return (
    <div>
      <Link to="/" className="back-link">
        &#8592; Back to dashboard
      </Link>

      <div className="page-header">
        <div className="page-header-top">
          <div>
            <div className="page-header-meta">
              <StatusBadge status={hc.Status} />
              {hc.Anonymous && (
                <span className="status-badge" style={{
                  background: 'var(--purple-dim)',
                  color: 'var(--purple)',
                }}>
                  <span>{'\uD83D\uDD12'}</span> Anonymous
                </span>
              )}
            </div>
            <h1>{hc.Name}</h1>
          </div>
          <div className={`page-header-score ${colorClass}`}>
            {total_votes > 0 ? average_score.toFixed(1) : '-'}
          </div>
        </div>
      </div>

      <div className="stats-row" style={{ marginBottom: '24px' }}>
        <div className="stat-item">
          <div className="stat-value">{participants}</div>
          <div className="stat-label">Participants</div>
        </div>
        <div className="stat-item">
          <div className="stat-value">{total_votes}</div>
          <div className="stat-label">Total Votes</div>
        </div>
        <div className="stat-item">
          <div className="stat-value">{results.length}</div>
          <div className="stat-label">Metrics</div>
        </div>
      </div>

      {!hc.Anonymous && participant_names && participant_names.length > 0 && (
        <div style={{ marginBottom: '24px' }}>
          <div className="section-title">Participants</div>
          <div className="participant-avatars">
            {participant_names.map((name) => (
              <div
                key={name}
                className="participant-avatar"
                style={{ backgroundColor: getAvatarColor(name) }}
                title={name}
              >
                {getInitial(name)}
              </div>
            ))}
          </div>
        </div>
      )}

      <div style={{ display: 'flex', gap: '12px', marginBottom: '24px', flexWrap: 'wrap' }}>
        {hc.Status === 'open' && (
          <Link to={`/healthcheck/${id}/vote`} className="btn btn-primary btn-lg">
            Cast Your Vote
          </Link>
        )}
        <a
          href={`/api/healthchecks/${id}/export`}
          className="btn btn-secondary"
          download
        >
          {'\u2B07'} Download CSV
        </a>
      </div>

      {(plugins || []).filter(p => p.nav_pos === 'healthcheck').length > 0 && (
        <div style={{ display: 'flex', gap: '12px', marginBottom: '24px', flexWrap: 'wrap' }}>
          {(plugins || [])
            .filter(p => p.nav_pos === 'healthcheck')
            .map(plugin => (
              <Link
                key={plugin.name}
                to={plugin.route.replace(':hcId', id || '')}
                className="btn btn-secondary"
              >
                {plugin.icon} {plugin.label}
              </Link>
            ))}
        </div>
      )}

      <VoteProgress results={results} totalMetrics={results.length} />

      {total_votes > 0 && results.length >= 3 && (
        <RadarChart results={results} />
      )}

      <div className="glass-card" style={{ padding: 0 }}>
        <div style={{ padding: '16px 20px 8px' }}>
          <div className="section-title" style={{ margin: 0 }}>Results</div>
        </div>
        <div style={{ padding: '0 20px 8px' }}>
          {results.map(metric => (
            <MetricRow key={metric.MetricName} metric={metric} />
          ))}
        </div>
      </div>

      {total_votes > 0 && id && (
        <DiscussionGuide healthcheckId={id} />
      )}

      {id && (
        <ActionItems
          healthcheckId={id}
          actions={data.actions || []}
          metricNames={results.map(r => r.MetricName)}
          onActionCreated={refetch}
        />
      )}
    </div>
  )
}
