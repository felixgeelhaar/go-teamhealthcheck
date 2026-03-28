import { useApi } from '../hooks/useApi'

interface Alert {
  metric: string
  severity: 'warning' | 'critical'
  message: string
  current_score: number
  trend: string
  delta: number
  predicted_score: number
}

interface AlertsResponse {
  alerts: Alert[]
  message?: string
}

interface Props {
  teamId: string
}

export function HealthAlerts({ teamId }: Props) {
  const { data } = useApi<AlertsResponse>(
    teamId ? `/api/teams/${teamId}/alerts` : null
  )

  if (!data || !data.alerts || data.alerts.length === 0) return null

  return (
    <div style={{ marginBottom: '24px' }}>
      <div className="section-title" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <span>{'\u26A0\uFE0F'}</span> Health Alerts
      </div>

      <div style={{ display: 'grid', gap: '8px' }}>
        {data.alerts.map((alert) => (
          <div
            key={alert.metric}
            className="glass-card"
            style={{
              padding: '12px 16px',
              borderLeft: `3px solid ${alert.severity === 'critical' ? 'var(--red)' : 'var(--yellow)'}`,
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: '12px' }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                  <span style={{
                    fontSize: '11px',
                    fontWeight: 700,
                    textTransform: 'uppercase' as const,
                    padding: '2px 8px',
                    borderRadius: '10px',
                    background: alert.severity === 'critical' ? 'rgba(239,68,68,0.15)' : 'rgba(234,179,8,0.15)',
                    color: alert.severity === 'critical' ? 'var(--red)' : 'var(--yellow)',
                  }}>
                    {alert.severity}
                  </span>
                  <span style={{ fontSize: '14px', fontWeight: 600 }}>{alert.metric}</span>
                </div>
                <div style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
                  {alert.message}
                </div>
              </div>
              <div style={{ textAlign: 'right', flexShrink: 0 }}>
                <div style={{
                  fontFamily: 'var(--font-mono)',
                  fontSize: '18px',
                  fontWeight: 700,
                  color: alert.current_score < 1.5 ? 'var(--red)' : alert.current_score < 2.5 ? 'var(--yellow)' : 'var(--green)',
                }}>
                  {alert.current_score.toFixed(1)}
                </div>
                <div style={{ fontSize: '11px', color: 'var(--text-tertiary)' }}>
                  → {alert.predicted_score.toFixed(1)}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
