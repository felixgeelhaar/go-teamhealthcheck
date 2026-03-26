import { useState, useCallback } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { useParticipant } from '../hooks/useParticipant'
import { MetricCard } from '../components/MetricCard'
import type { HealthCheckDetail, VotePayload } from '../types'

type VoteColor = 'green' | 'yellow' | 'red'

interface VoteState {
  color: VoteColor | null
  comment: string
}

export function VotingPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { name: participant } = useParticipant()
  const { data, loading } = useApi<HealthCheckDetail>(
    id ? `/api/healthchecks/${id}` : null
  )

  const [votes, setVotes] = useState<Record<string, VoteState>>({})
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)

  const setVoteColor = useCallback((metricName: string, color: VoteColor) => {
    setVotes(prev => ({
      ...prev,
      [metricName]: { ...prev[metricName], color, comment: prev[metricName]?.comment || '' },
    }))
  }, [])

  const setVoteComment = useCallback((metricName: string, comment: string) => {
    setVotes(prev => ({
      ...prev,
      [metricName]: { ...prev[metricName], color: prev[metricName]?.color || null, comment },
    }))
  }, [])

  const handleSubmit = async () => {
    if (!id || !participant) return

    const votesToSubmit = Object.entries(votes).filter(
      ([, v]) => v.color !== null
    ) as [string, { color: VoteColor; comment: string }][]

    if (votesToSubmit.length === 0) return

    setSubmitting(true)
    setSubmitError(null)

    try {
      for (const [metricName, vote] of votesToSubmit) {
        const payload: VotePayload = {
          participant,
          metric_name: metricName,
          color: vote.color,
        }
        if (vote.comment.trim()) {
          payload.comment = vote.comment.trim()
        }

        const res = await fetch(`/api/healthchecks/${id}/vote`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        })

        if (!res.ok) {
          const text = await res.text()
          throw new Error(text || `Failed to submit vote for ${metricName}`)
        }
      }

      navigate(`/healthcheck/${id}`)
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'Failed to submit votes')
    } finally {
      setSubmitting(false)
    }
  }

  if (loading || !data) {
    return (
      <div className="loading">
        <div className="loading-spinner" />
        Loading...
      </div>
    )
  }

  const { healthcheck: hc, template } = data
  const metrics = template.Metrics || []
  const totalMetrics = metrics.length
  const votedCount = Object.values(votes).filter(v => v.color !== null).length
  const allVoted = votedCount === totalMetrics && totalMetrics > 0

  return (
    <div>
      <Link to={`/healthcheck/${id}`} className="back-link">
        &#8592; Back to results
      </Link>

      <div className="page-header">
        <h1>Vote: {hc.Name}</h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '14px', marginTop: '4px' }}>
          Voting as <strong style={{ color: 'var(--text-primary)' }}>{participant}</strong> &middot; {votedCount} of {totalMetrics} metrics
        </p>
      </div>

      {metrics.map((metric) => (
        <MetricCard
          key={metric.Name}
          metric={metric}
          selectedColor={votes[metric.Name]?.color || null}
          onColorSelect={(color) => setVoteColor(metric.Name, color)}
          comment={votes[metric.Name]?.comment || ''}
          onCommentChange={(comment) => setVoteComment(metric.Name, comment)}
        />
      ))}

      {submitError && (
        <div style={{
          marginTop: '16px',
          padding: '12px 16px',
          background: 'var(--red-dim)',
          border: '1px solid rgba(239, 68, 68, 0.2)',
          borderRadius: 'var(--radius-md)',
          color: 'var(--red)',
          fontSize: '14px',
        }}>
          {submitError}
        </div>
      )}

      <div style={{ marginTop: '24px', textAlign: 'center' }}>
        <button
          className="btn btn-primary btn-lg"
          onClick={handleSubmit}
          disabled={!allVoted || submitting}
        >
          {submitting ? 'Submitting...' : `Submit All Votes (${votedCount}/${totalMetrics})`}
        </button>
      </div>
    </div>
  )
}
