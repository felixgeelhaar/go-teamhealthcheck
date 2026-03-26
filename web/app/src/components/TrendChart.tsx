import { scoreColorClass } from '../utils'

interface SessionSummary {
  id: string
  name: string
  date: string
  avg_score: number
  voters: number
  status: string
}

interface MetricTrend {
  MetricName: string
  Sessions: { HealthCheckID: string; HealthCheckName: string; Score: number; Date: string }[]
  Tendency: 'improving' | 'stable' | 'declining'
  Delta: number
}

interface Props {
  sessions: SessionSummary[]
  trends: MetricTrend[]
}

const tendencyIcon = (t: string) => {
  switch (t) {
    case 'improving': return '↑'
    case 'declining': return '↓'
    default: return '→'
  }
}

const tendencyColor = (t: string) => {
  switch (t) {
    case 'improving': return 'var(--green)'
    case 'declining': return 'var(--red)'
    default: return 'var(--text-secondary)'
  }
}

export function TrendChart({ sessions, trends }: Props) {
  if (sessions.length === 0) return null

  const maxScore = 3.0

  return (
    <div style={{ marginBottom: '32px' }}>
      <div className="section-title">Team Health Over Time</div>

      {/* Overall score trend bar chart */}
      <div className="glass-card" style={{ marginBottom: '16px' }}>
        <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '12px' }}>
          Average Score per Session
        </div>
        <div style={{ display: 'flex', gap: '8px', alignItems: 'flex-end', minHeight: '120px' }}>
          {sessions.map((s) => {
            const heightPct = s.avg_score > 0 ? (s.avg_score / maxScore) * 100 : 2
            const colorClass = scoreColorClass(s.avg_score)
            const barColor = colorClass === 'score-green' ? 'var(--green)'
              : colorClass === 'score-yellow' ? 'var(--yellow)' : 'var(--red)'
            return (
              <div key={s.id} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '4px' }}>
                <div style={{
                  fontSize: '11px',
                  fontFamily: 'var(--font-mono)',
                  fontWeight: 700,
                  color: barColor,
                }}>
                  {s.avg_score > 0 ? s.avg_score.toFixed(1) : '-'}
                </div>
                <div style={{
                  width: '100%',
                  maxWidth: '48px',
                  height: `${heightPct}px`,
                  minHeight: '4px',
                  background: barColor,
                  borderRadius: '4px 4px 2px 2px',
                  opacity: 0.8,
                  transition: 'height 0.3s ease',
                }} />
                <div style={{
                  fontSize: '10px',
                  color: 'var(--text-secondary)',
                  textAlign: 'center',
                  lineHeight: '1.2',
                  maxWidth: '60px',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                }}>
                  {s.name}
                </div>
              </div>
            )
          })}
        </div>
      </div>

      {/* Per-metric trend indicators */}
      <div className="glass-card">
        <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '12px' }}>
          Metric Trends
        </div>
        <div style={{ display: 'grid', gap: '8px' }}>
          {trends.map((t) => {
            const latestScore = t.Sessions.length > 0 ? t.Sessions[t.Sessions.length - 1].Score : 0
            const colorClass = scoreColorClass(latestScore)
            return (
              <div key={t.MetricName} style={{
                display: 'grid',
                gridTemplateColumns: '1fr auto auto auto',
                gap: '12px',
                alignItems: 'center',
                padding: '6px 0',
                borderBottom: '1px solid rgba(255,255,255,0.04)',
              }}>
                <div style={{ fontSize: '13px', fontWeight: 500 }}>{t.MetricName}</div>

                {/* Mini sparkline - dots for each session */}
                <div style={{ display: 'flex', gap: '3px', alignItems: 'center' }}>
                  {t.Sessions.map((s, i) => {
                    const size = 6
                    const sc = scoreColorClass(s.Score)
                    const c = sc === 'score-green' ? 'var(--green)'
                      : sc === 'score-yellow' ? 'var(--yellow)' : 'var(--red)'
                    return (
                      <div key={i} style={{
                        width: `${size}px`,
                        height: `${size}px`,
                        borderRadius: '50%',
                        background: s.Score > 0 ? c : 'var(--text-secondary)',
                        opacity: 0.8,
                      }} title={`${s.HealthCheckName}: ${s.Score.toFixed(1)}`} />
                    )
                  })}
                </div>

                <div className={`metric-score ${colorClass}`} style={{ fontSize: '14px', minWidth: '32px', textAlign: 'right' }}>
                  {latestScore > 0 ? latestScore.toFixed(1) : '-'}
                </div>

                <div style={{
                  fontSize: '14px',
                  fontWeight: 700,
                  color: tendencyColor(t.Tendency),
                  minWidth: '20px',
                  textAlign: 'center',
                }} title={`${t.Tendency} (${t.Delta > 0 ? '+' : ''}${t.Delta.toFixed(1)})`}>
                  {tendencyIcon(t.Tendency)}
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
