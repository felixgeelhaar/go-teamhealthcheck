import type { MetricResult } from '../types'

interface Props {
  results: MetricResult[]
  totalMetrics: number
}

export function VoteProgress({ results, totalMetrics }: Props) {
  const participants = new Set<string>()
  const votedMetrics = results.filter(r => r.TotalVotes > 0).length

  // Collect unique participants from comments (rough proxy)
  // The actual participant count comes from the API response
  results.forEach(r => {
    if (r.TotalVotes > 0) participants.add(r.MetricName)
  })

  const pct = totalMetrics > 0 ? (votedMetrics / totalMetrics) * 100 : 0

  return (
    <div style={{
      padding: '16px',
      backgroundColor: '#f9fafb',
      borderRadius: '8px',
      marginBottom: '24px',
    }}>
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        marginBottom: '8px',
        fontSize: '14px',
        color: '#374151',
      }}>
        <span>Metrics with votes</span>
        <span>{votedMetrics} / {totalMetrics}</span>
      </div>
      <div style={{
        height: '8px',
        backgroundColor: '#e5e7eb',
        borderRadius: '4px',
        overflow: 'hidden',
      }}>
        <div style={{
          width: `${pct}%`,
          height: '100%',
          backgroundColor: pct === 100 ? '#22c55e' : '#3b82f6',
          transition: 'width 0.3s',
          borderRadius: '4px',
        }} />
      </div>
    </div>
  )
}
