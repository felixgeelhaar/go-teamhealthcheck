import type { MetricResult } from '../types'

interface Props {
  metric: MetricResult
}

export function MetricRow({ metric }: Props) {
  const total = metric.TotalVotes || 1
  const greenPct = (metric.GreenCount / total) * 100
  const yellowPct = (metric.YellowCount / total) * 100
  const redPct = (metric.RedCount / total) * 100

  const scoreColor =
    metric.Score >= 2.5 ? '#22c55e' :
    metric.Score >= 1.5 ? '#eab308' : '#ef4444'

  return (
    <div style={{
      display: 'grid',
      gridTemplateColumns: '200px 1fr 60px',
      gap: '12px',
      alignItems: 'center',
      padding: '8px 0',
      borderBottom: '1px solid #e5e7eb',
    }}>
      <div>
        <div style={{ fontWeight: 600, fontSize: '14px' }}>{metric.MetricName}</div>
        <div style={{ fontSize: '11px', color: '#6b7280', marginTop: '2px' }}>
          {metric.TotalVotes} vote{metric.TotalVotes !== 1 ? 's' : ''}
        </div>
      </div>

      <div style={{
        display: 'flex',
        height: '24px',
        borderRadius: '4px',
        overflow: 'hidden',
        backgroundColor: '#f3f4f6',
      }}>
        {metric.GreenCount > 0 && (
          <div style={{ width: `${greenPct}%`, backgroundColor: '#22c55e', transition: 'width 0.3s' }} />
        )}
        {metric.YellowCount > 0 && (
          <div style={{ width: `${yellowPct}%`, backgroundColor: '#eab308', transition: 'width 0.3s' }} />
        )}
        {metric.RedCount > 0 && (
          <div style={{ width: `${redPct}%`, backgroundColor: '#ef4444', transition: 'width 0.3s' }} />
        )}
      </div>

      <div style={{
        fontWeight: 700,
        fontSize: '16px',
        textAlign: 'center',
        color: scoreColor,
      }}>
        {metric.TotalVotes > 0 ? metric.Score.toFixed(1) : '-'}
      </div>
    </div>
  )
}
