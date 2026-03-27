import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import type { Template, CreateHealthCheckPayload } from '../types'

export function CreateHealthCheck() {
  const { teamId } = useParams<{ teamId: string }>()
  const navigate = useNavigate()
  const { data: templates, loading: templatesLoading } = useApi<Template[]>('/api/templates')

  const [name, setName] = useState('')
  const [templateId, setTemplateId] = useState('')
  const [anonymous, setAnonymous] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const selectedTemplate = templates?.find(t => t.ID === templateId)

  const isValid = () => {
    return name.trim().length > 0 && templateId.length > 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!isValid() || !teamId) return

    setSubmitting(true)
    setError(null)

    const payload: CreateHealthCheckPayload = {
      name: name.trim(),
      template_id: templateId,
      anonymous,
    }

    try {
      const res = await fetch(`/api/teams/${teamId}/healthchecks`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })

      if (res.ok) {
        const result = await res.json()
        if (result && result.ID) {
          navigate(`/healthcheck/${result.ID}`)
        } else {
          navigate('/')
        }
      } else {
        const text = await res.text()
        setError(text || 'Failed to create health check')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Network error')
    } finally {
      setSubmitting(false)
    }
  }

  if (templatesLoading) {
    return (
      <div className="loading">
        <div className="loading-spinner" />
        Loading...
      </div>
    )
  }

  return (
    <div>
      <Link to="/" className="back-link">
        &#8592; Back to dashboard
      </Link>

      <div className="page-header">
        <h1>New Health Check</h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '14px', marginTop: '4px' }}>
          Start a new health check session for your team.
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="glass-card">
          <div className="form-group">
            <label className="form-label">Health Check Name</label>
            <input
              type="text"
              className="form-input"
              placeholder="e.g., Sprint 42 Review"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>

          <div className="form-group" style={{ marginBottom: 0 }}>
            <label className="form-label">Template</label>
            <select
              className="form-select"
              value={templateId}
              onChange={(e) => setTemplateId(e.target.value)}
            >
              <option value="">Select a template...</option>
              {templates?.map(t => (
                <option key={t.ID} value={t.ID}>
                  {t.Name}{t.BuiltIn ? ' (built-in)' : ''}
                </option>
              ))}
            </select>
          </div>

          {selectedTemplate && selectedTemplate.Metrics && (
            <div className="template-preview">
              <div className="template-preview-title">
                Template Metrics ({selectedTemplate.Metrics.length})
              </div>
              {selectedTemplate.Metrics.map((m) => (
                <div key={m.Name} className="template-preview-item">
                  {m.Name}
                </div>
              ))}
            </div>
          )}

          <div style={{ marginTop: '20px', paddingTop: '20px', borderTop: '1px solid var(--glass-border)' }}>
            <label
              style={{
                display: 'flex',
                alignItems: 'flex-start',
                gap: '12px',
                cursor: 'pointer',
              }}
            >
              <span
                role="checkbox"
                aria-checked={anonymous}
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === ' ' || e.key === 'Enter') {
                    e.preventDefault()
                    setAnonymous(!anonymous)
                  }
                }}
                onClick={() => setAnonymous(!anonymous)}
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '20px',
                  height: '20px',
                  borderRadius: 'var(--radius-sm)',
                  border: anonymous
                    ? '2px solid var(--blue)'
                    : '2px solid var(--glass-border)',
                  background: anonymous
                    ? 'var(--blue)'
                    : 'rgba(255, 255, 255, 0.03)',
                  transition: 'all var(--transition-fast)',
                  flexShrink: 0,
                  marginTop: '2px',
                  color: 'white',
                  fontSize: '12px',
                  fontWeight: 700,
                }}
              >
                {anonymous ? '\u2713' : ''}
              </span>
              <span>
                <span style={{ display: 'block', fontSize: '14px', fontWeight: 600, color: 'var(--text-primary)' }}>
                  Anonymous voting
                </span>
                <span style={{ display: 'block', fontSize: '13px', color: 'var(--text-tertiary)', marginTop: '2px' }}>
                  Participant names will be hidden in results
                </span>
              </span>
            </label>
          </div>
        </div>

        {error && (
          <div style={{
            marginTop: '16px',
            padding: '12px 16px',
            background: 'var(--red-dim)',
            border: '1px solid rgba(239, 68, 68, 0.2)',
            borderRadius: 'var(--radius-md)',
            color: 'var(--red)',
            fontSize: '14px',
          }}>
            {error}
          </div>
        )}

        <div style={{ marginTop: '24px' }}>
          <button
            type="submit"
            className="btn btn-primary btn-lg"
            disabled={!isValid() || submitting}
          >
            {submitting ? 'Creating...' : 'Create Health Check'}
          </button>
        </div>
      </form>
    </div>
  )
}
