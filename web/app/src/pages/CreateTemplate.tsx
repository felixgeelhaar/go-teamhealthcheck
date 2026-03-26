import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import type { CreateTemplateMetric, CreateTemplatePayload } from '../types'

interface MetricFormState {
  name: string
  description_good: string
  description_bad: string
}

const emptyMetric = (): MetricFormState => ({
  name: '',
  description_good: '',
  description_bad: '',
})

export function CreateTemplate() {
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [metrics, setMetrics] = useState<MetricFormState[]>([emptyMetric()])
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const addMetric = () => {
    setMetrics([...metrics, emptyMetric()])
  }

  const removeMetric = (index: number) => {
    if (metrics.length <= 1) return
    setMetrics(metrics.filter((_, i) => i !== index))
  }

  const updateMetric = (index: number, field: keyof MetricFormState, value: string) => {
    setMetrics(metrics.map((m, i) => i === index ? { ...m, [field]: value } : m))
  }

  const isValid = () => {
    if (!name.trim()) return false
    return metrics.every(m => m.name.trim() && m.description_good.trim() && m.description_bad.trim())
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!isValid()) return

    setSubmitting(true)
    setError(null)

    const payload: CreateTemplatePayload = {
      name: name.trim(),
      description: description.trim(),
      metrics: metrics.map((m): CreateTemplateMetric => ({
        name: m.name.trim(),
        description_good: m.description_good.trim(),
        description_bad: m.description_bad.trim(),
      })),
    }

    try {
      const res = await fetch('/api/templates', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })

      if (res.ok) {
        navigate('/')
      } else {
        const text = await res.text()
        setError(text || 'Failed to create template')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Network error')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div>
      <Link to="/" className="back-link">
        &#8592; Back to dashboard
      </Link>

      <div className="page-header">
        <h1>Create Template</h1>
        <p style={{ color: 'var(--text-secondary)', fontSize: '14px', marginTop: '4px' }}>
          Define the metrics your team will be evaluated on.
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="glass-card" style={{ marginBottom: '24px' }}>
          <div className="form-group">
            <label className="form-label">Template Name</label>
            <input
              type="text"
              className="form-input"
              placeholder="e.g., Sprint Health Check"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          <div className="form-group" style={{ marginBottom: 0 }}>
            <label className="form-label">Description</label>
            <textarea
              className="form-textarea"
              placeholder="What is this template for?"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={2}
            />
          </div>
        </div>

        <div className="section-title">Metrics</div>

        {metrics.map((metric, index) => (
          <div key={index} className="metric-form-item">
            <div className="metric-form-header">
              <span className="metric-form-number">Metric {index + 1}</span>
              {metrics.length > 1 && (
                <button
                  type="button"
                  className="btn btn-danger btn-sm"
                  onClick={() => removeMetric(index)}
                >
                  Remove
                </button>
              )}
            </div>

            <div className="form-group">
              <label className="form-label">Metric Name</label>
              <input
                type="text"
                className="form-input"
                placeholder="e.g., Team Collaboration"
                value={metric.name}
                onChange={(e) => updateMetric(index, 'name', e.target.value)}
              />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
              <div className="form-group" style={{ marginBottom: 0 }}>
                <label className="form-label" style={{ color: 'var(--green)' }}>
                  What Good Looks Like
                </label>
                <textarea
                  className="form-textarea"
                  placeholder="Describe the ideal state..."
                  value={metric.description_good}
                  onChange={(e) => updateMetric(index, 'description_good', e.target.value)}
                  rows={2}
                />
              </div>
              <div className="form-group" style={{ marginBottom: 0 }}>
                <label className="form-label" style={{ color: 'var(--red)' }}>
                  What Bad Looks Like
                </label>
                <textarea
                  className="form-textarea"
                  placeholder="Describe problematic state..."
                  value={metric.description_bad}
                  onChange={(e) => updateMetric(index, 'description_bad', e.target.value)}
                  rows={2}
                />
              </div>
            </div>
          </div>
        ))}

        <button
          type="button"
          className="btn btn-secondary"
          onClick={addMetric}
          style={{ marginTop: '4px', marginBottom: '24px' }}
        >
          + Add Metric
        </button>

        {error && (
          <div style={{
            marginBottom: '16px',
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

        <button
          type="submit"
          className="btn btn-primary btn-lg"
          disabled={!isValid() || submitting}
        >
          {submitting ? 'Creating...' : 'Create Template'}
        </button>
      </form>
    </div>
  )
}
