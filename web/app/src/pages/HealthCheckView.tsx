import { useParams, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { useWebSocket } from '../hooks/useWebSocket'
import { MetricRow } from '../components/MetricRow'
import { VoteProgress } from '../components/VoteProgress'
import type { HealthCheckResults, WSEvent } from '../types'

export function HealthCheckView() {
  const { id } = useParams<{ id: string }>()
  const { data, refetch } = useApi<HealthCheckResults>(
    id ? `/api/healthchecks/${id}/results` : null
  )

  useWebSocket((event: WSEvent) => {
    if (event.healthcheck_id === id) {
      refetch()
    }
  })

  if (!data) {
    return <div style={{ color: '#9ca3af', padding: '24px' }}>Loading...</div>
  }

  const { healthcheck: hc, results, average_score, participants, total_votes } = data

  const scoreColor =
    average_score >= 2.5 ? '#22c55e' :
    average_score >= 1.5 ? '#eab308' : '#ef4444'

  return (
    <div>
      <Link to="/" style={{ color: '#3b82f6', fontSize: '14px', textDecoration: 'none' }}>
        &larr; Back
      </Link>

      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'flex-start',
        marginTop: '16px',
        marginBottom: '24px',
      }}>
        <div>
          <h1 style={{ fontSize: '24px', fontWeight: 700, margin: 0 }}>{hc.Name}</h1>
          <div style={{ fontSize: '14px', color: '#6b7280', marginTop: '4px' }}>
            {hc.Status} &middot; {participants} participant{participants !== 1 ? 's' : ''} &middot; {total_votes} total votes
          </div>
        </div>
        <div style={{
          fontSize: '32px',
          fontWeight: 800,
          color: scoreColor,
        }}>
          {total_votes > 0 ? average_score.toFixed(1) : '-'}
        </div>
      </div>

      <VoteProgress
        results={results}
        totalMetrics={results.length}
      />

      <div style={{ marginBottom: '8px' }}>
        <div style={{
          display: 'grid',
          gridTemplateColumns: '200px 1fr 60px',
          gap: '12px',
          padding: '0 0 8px',
          fontSize: '12px',
          fontWeight: 600,
          color: '#9ca3af',
          textTransform: 'uppercase',
          letterSpacing: '0.05em',
        }}>
          <div>Metric</div>
          <div>Votes</div>
          <div style={{ textAlign: 'center' }}>Score</div>
        </div>
        {results.map(metric => (
          <MetricRow key={metric.MetricName} metric={metric} />
        ))}
      </div>
    </div>
  )
}
