import { useApi } from '../hooks/useApi'
import type { DiscussionGuideResponse } from '../types'

interface Props {
  healthcheckId: string
}

function reasonColor(reason: string): { background: string; color: string; border: string } {
  const r = reason.toLowerCase()
  if (r.includes('low score') && r.includes('disagreement')) {
    return {
      background: 'rgba(249, 115, 22, 0.15)',
      color: '#f97316',
      border: '1px solid rgba(249, 115, 22, 0.25)',
    }
  }
  if (r.includes('disagreement')) {
    return {
      background: 'var(--yellow-dim)',
      color: 'var(--yellow)',
      border: '1px solid rgba(234, 179, 8, 0.25)',
    }
  }
  return {
    background: 'var(--red-dim)',
    color: 'var(--red)',
    border: '1px solid rgba(239, 68, 68, 0.25)',
  }
}

function scoreColor(score: number): string {
  if (score >= 2.5) return 'var(--green)'
  if (score >= 1.5) return 'var(--yellow)'
  return 'var(--red)'
}

export function DiscussionGuide({ healthcheckId }: Props) {
  const { data, loading, error } = useApi<DiscussionGuideResponse>(
    `/api/healthchecks/${healthcheckId}/discussion`
  )

  if (loading) {
    return (
      <div className="glass-card" style={{ marginTop: '24px' }}>
        <div className="section-title" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span>{'\u2728'}</span> Discussion Guide
        </div>
        <div className="loading" style={{ padding: '24px 0' }}>
          <div className="loading-spinner" />
          Generating discussion topics...
        </div>
      </div>
    )
  }

  if (error || !data) {
    return null
  }

  const topics = data.topics || []

  return (
    <div className="glass-card" style={{ marginTop: '24px' }}>
      <div className="section-title" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <span>{'\u2728'}</span> Discussion Guide
      </div>

      {topics.length === 0 ? (
        <div className="empty-state" style={{ padding: '32px 20px' }}>
          <div className="empty-state-icon">{'\u2705'}</div>
          <div className="empty-state-text">
            No concerns detected — your team is doing great!
          </div>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {topics.map((topic) => {
            const pillStyle = reasonColor(topic.reason)
            return (
              <div
                key={topic.priority}
                style={{
                  padding: '16px',
                  background: 'rgba(255, 255, 255, 0.02)',
                  border: '1px solid var(--glass-border)',
                  borderRadius: 'var(--radius-md)',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '8px', flexWrap: 'wrap' }}>
                  <span
                    style={{
                      display: 'inline-flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      width: '24px',
                      height: '24px',
                      borderRadius: 'var(--radius-full)',
                      background: 'var(--blue-dim)',
                      color: 'var(--blue)',
                      fontSize: '12px',
                      fontWeight: 700,
                      fontFamily: 'var(--font-mono)',
                      flexShrink: 0,
                    }}
                  >
                    {topic.priority}
                  </span>
                  <span style={{ fontWeight: 700, fontSize: '15px', color: 'var(--text-primary)' }}>
                    {topic.metric}
                  </span>
                  <span
                    style={{
                      display: 'inline-flex',
                      padding: '2px 10px',
                      borderRadius: 'var(--radius-full)',
                      fontSize: '11px',
                      fontWeight: 600,
                      textTransform: 'uppercase' as const,
                      letterSpacing: '0.05em',
                      background: pillStyle.background,
                      color: pillStyle.color,
                      border: pillStyle.border,
                    }}
                  >
                    {topic.reason}
                  </span>
                  <span
                    style={{
                      fontFamily: 'var(--font-mono)',
                      fontWeight: 800,
                      fontSize: '14px',
                      color: scoreColor(topic.score),
                    }}
                  >
                    {topic.score.toFixed(1)}
                  </span>
                </div>

                {topic.data_points && topic.data_points.length > 0 && (
                  <div style={{ marginBottom: '8px' }}>
                    {topic.data_points.map((dp, i) => (
                      <div
                        key={i}
                        style={{
                          fontSize: '12px',
                          color: 'var(--text-tertiary)',
                          fontFamily: 'var(--font-mono)',
                        }}
                      >
                        {dp}
                      </div>
                    ))}
                  </div>
                )}

                {topic.suggested_questions && topic.suggested_questions.length > 0 && (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
                    {topic.suggested_questions.map((q, i) => (
                      <div
                        key={i}
                        style={{
                          fontSize: '13px',
                          color: 'var(--text-secondary)',
                          fontStyle: 'italic',
                          display: 'flex',
                          alignItems: 'flex-start',
                          gap: '6px',
                        }}
                      >
                        <span style={{ color: 'var(--text-tertiary)', flexShrink: 0 }}>?</span>
                        <span>{q}</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
