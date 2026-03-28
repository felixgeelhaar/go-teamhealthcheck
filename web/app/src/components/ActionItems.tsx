import { useState } from 'react'
import { useParticipant } from '../hooks/useParticipant'
import { getAvatarColor, getInitial, timeAgo } from '../utils'
import type { Action } from '../types'

interface Props {
  healthcheckId: string
  actions: Action[]
  metricNames: string[]
  onActionCreated: () => void
}

export function ActionItems({ healthcheckId, actions, metricNames, onActionCreated }: Props) {
  const { name } = useParticipant()
  const [description, setDescription] = useState('')
  const [assignee, setAssignee] = useState(name)
  const [metricName, setMetricName] = useState(metricNames[0] || '')
  const [submitting, setSubmitting] = useState(false)
  const [completing, setCompleting] = useState<string | null>(null)
  const [generating, setGenerating] = useState(false)

  const handleAdd = async () => {
    if (!description.trim() || !assignee.trim() || !metricName) return
    setSubmitting(true)
    try {
      const res = await fetch(`/api/healthchecks/${healthcheckId}/actions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          metric_name: metricName,
          description: description.trim(),
          assignee: assignee.trim(),
        }),
      })
      if (res.ok) {
        setDescription('')
        onActionCreated()
      }
    } finally {
      setSubmitting(false)
    }
  }

  const handleComplete = async (actionId: string) => {
    setCompleting(actionId)
    try {
      const res = await fetch(`/api/actions/${actionId}/complete`, {
        method: 'POST',
      })
      if (res.ok) {
        onActionCreated()
      }
    } finally {
      setCompleting(null)
    }
  }

  const handleGenerate = async () => {
    setGenerating(true)
    try {
      const res = await fetch(`/api/healthchecks/${healthcheckId}/generate-actions`, {
        method: 'POST',
      })
      if (res.ok) {
        onActionCreated()
      }
    } finally {
      setGenerating(false)
    }
  }

  const pendingActions = (actions || []).filter(a => !a.Completed)
  const completedActions = (actions || []).filter(a => a.Completed)

  return (
    <div className="glass-card" style={{ marginTop: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
        <div className="section-title" style={{ display: 'flex', alignItems: 'center', gap: '8px', margin: 0 }}>
          <span>{'\uD83D\uDCCB'}</span> Action Items
        </div>
        <button
          className="btn btn-secondary btn-sm"
          onClick={handleGenerate}
          disabled={generating}
          style={{ whiteSpace: 'nowrap' }}
        >
          {generating ? 'Generating...' : '\u2728 Generate Actions'}
        </button>
      </div>

      {pendingActions.length === 0 && completedActions.length === 0 && (
        <div style={{ color: 'var(--text-tertiary)', fontSize: '14px', marginBottom: '16px' }}>
          No action items yet. Add one below to track follow-ups.
        </div>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
        {pendingActions.map(action => (
          <ActionRow
            key={action.ID}
            action={action}
            completing={completing === action.ID}
            onComplete={() => handleComplete(action.ID)}
          />
        ))}
        {completedActions.map(action => (
          <ActionRow
            key={action.ID}
            action={action}
            completing={completing === action.ID}
            onComplete={() => handleComplete(action.ID)}
          />
        ))}
      </div>

      <div
        style={{
          marginTop: '16px',
          paddingTop: '16px',
          borderTop: '1px solid var(--glass-border)',
          display: 'flex',
          gap: '8px',
          flexWrap: 'wrap',
          alignItems: 'flex-end',
        }}
      >
        <div style={{ flex: '2 1 180px', minWidth: 0 }}>
          <label className="form-label">Description</label>
          <input
            className="form-input"
            type="text"
            placeholder="What needs to be done?"
            value={description}
            onChange={e => setDescription(e.target.value)}
            onKeyDown={e => { if (e.key === 'Enter') handleAdd() }}
          />
        </div>
        <div style={{ flex: '1 1 120px', minWidth: 0 }}>
          <label className="form-label">Metric</label>
          <select
            className="form-select"
            value={metricName}
            onChange={e => setMetricName(e.target.value)}
          >
            {metricNames.map(m => (
              <option key={m} value={m}>{m}</option>
            ))}
          </select>
        </div>
        <div style={{ flex: '1 1 120px', minWidth: 0 }}>
          <label className="form-label">Assignee</label>
          <input
            className="form-input"
            type="text"
            placeholder="Who?"
            value={assignee}
            onChange={e => setAssignee(e.target.value)}
            onKeyDown={e => { if (e.key === 'Enter') handleAdd() }}
          />
        </div>
        <button
          className="btn btn-primary btn-sm"
          onClick={handleAdd}
          disabled={submitting || !description.trim() || !assignee.trim()}
          style={{ flexShrink: 0, alignSelf: 'flex-end' }}
        >
          {submitting ? 'Adding...' : 'Add'}
        </button>
      </div>
    </div>
  )
}

function ActionRow({
  action,
  completing,
  onComplete,
}: {
  action: Action
  completing: boolean
  onComplete: () => void
}) {
  const completed = action.Completed

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: '10px',
        padding: '10px 12px',
        background: 'rgba(255, 255, 255, 0.02)',
        border: '1px solid var(--glass-border)',
        borderRadius: 'var(--radius-md)',
        opacity: completed ? 0.5 : 1,
        transition: 'opacity var(--transition-fast)',
      }}
    >
      <button
        onClick={onComplete}
        disabled={completing || completed}
        style={{
          width: '20px',
          height: '20px',
          borderRadius: 'var(--radius-sm)',
          border: completed
            ? '2px solid var(--green)'
            : '2px solid var(--glass-border)',
          background: completed ? 'var(--green)' : 'transparent',
          cursor: completed ? 'default' : 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          flexShrink: 0,
          transition: 'all var(--transition-fast)',
          color: 'white',
          fontSize: '12px',
          fontWeight: 700,
          padding: 0,
        }}
        aria-label={completed ? 'Completed' : 'Mark as complete'}
      >
        {completed ? '\u2713' : ''}
      </button>

      <span
        style={{
          flex: 1,
          fontSize: '14px',
          color: completed ? 'var(--text-muted)' : 'var(--text-primary)',
          textDecoration: completed ? 'line-through' : 'none',
          minWidth: 0,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
        }}
      >
        {action.Description}
      </span>

      <span
        style={{
          display: 'inline-flex',
          padding: '2px 8px',
          borderRadius: 'var(--radius-full)',
          fontSize: '11px',
          fontWeight: 600,
          background: 'var(--blue-dim)',
          color: 'var(--blue)',
          border: '1px solid rgba(59, 130, 246, 0.2)',
          whiteSpace: 'nowrap',
          flexShrink: 0,
        }}
      >
        {action.MetricName}
      </span>

      <div
        style={{
          width: '24px',
          height: '24px',
          borderRadius: '50%',
          backgroundColor: getAvatarColor(action.Assignee),
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '11px',
          fontWeight: 700,
          color: 'white',
          flexShrink: 0,
        }}
        title={action.Assignee}
      >
        {getInitial(action.Assignee)}
      </div>

      <span
        style={{
          fontSize: '12px',
          color: 'var(--text-tertiary)',
          whiteSpace: 'nowrap',
          flexShrink: 0,
          fontFamily: 'var(--font-mono)',
        }}
      >
        {timeAgo(action.CreatedAt)}
      </span>
    </div>
  )
}
