import type { MetricResult } from '../types'

interface Props {
  results: MetricResult[]
  totalMetrics: number
}

export function VoteProgress({ results, totalMetrics }: Props) {
  const votedMetrics = results.filter(r => r.TotalVotes > 0).length
  const pct = totalMetrics > 0 ? (votedMetrics / totalMetrics) * 100 : 0
  const isComplete = pct === 100

  return (
    <div className="glass-card vote-progress-bar">
      <div className="vote-progress-header">
        <span>Metrics with votes</span>
        <span>{votedMetrics} / {totalMetrics}</span>
      </div>
      <div className="vote-progress-track">
        <div
          className="vote-progress-fill"
          style={{
            width: `${pct}%`,
            background: isComplete ? 'var(--green)' : 'var(--blue)',
          }}
        />
      </div>
    </div>
  )
}
