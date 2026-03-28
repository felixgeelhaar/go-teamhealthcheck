import { useState, useRef } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import { useParticipant } from '../hooks/useParticipant'
import { useWebSocket } from '../hooks/useWebSocket'
import { getAvatarColor, getInitial } from '../utils'
import type { RetroResponse, RetroItem, HealthCheckResults, WSEvent } from '../types'

const CATEGORIES = [
  { key: 'went_well' as const, label: 'Went Well', color: 'var(--green)', dim: 'var(--green-dim)', glow: 'var(--green-glow)' },
  { key: 'to_improve' as const, label: 'To Improve', color: 'var(--yellow)', dim: 'var(--yellow-dim)', glow: 'var(--yellow-glow)' },
  { key: 'action_item' as const, label: 'Action Items', color: 'var(--blue)', dim: 'var(--blue-dim)', glow: 'var(--blue-glow)' },
] as const

export function RetroPage() {
  const { hcId } = useParams<{ hcId: string }>()
  const { name } = useParticipant()
  const [newItemText, setNewItemText] = useState<Record<string, string>>({
    went_well: '',
    to_improve: '',
    action_item: '',
  })
  const [submitting, setSubmitting] = useState<Record<string, boolean>>({})
  const formRefs = useRef<Record<string, HTMLInputElement | null>>({})

  const { data: retroData, loading: retroLoading, refetch } = useApi<RetroResponse>(
    hcId ? `/api/retro/${hcId}` : null
  )
  const { data: hcData } = useApi<HealthCheckResults>(
    hcId ? `/api/healthchecks/${hcId}/results` : null
  )

  useWebSocket((event: WSEvent) => {
    if (event.healthcheck_id === hcId) {
      refetch()
    }
  })

  const handleStartRetro = async () => {
    if (!hcId) return
    setSubmitting(prev => ({ ...prev, start: true }))
    try {
      const res = await fetch(`/api/retro/${hcId}`, { method: 'POST' })
      if (res.ok) {
        refetch()
      }
    } finally {
      setSubmitting(prev => ({ ...prev, start: false }))
    }
  }

  const handleAddItem = async (category: RetroItem['category']) => {
    const text = newItemText[category]?.trim()
    if (!text || !hcId) return
    setSubmitting(prev => ({ ...prev, [category]: true }))
    try {
      const res = await fetch(`/api/retro/${hcId}/items`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ category, text, author: name }),
      })
      if (res.ok) {
        setNewItemText(prev => ({ ...prev, [category]: '' }))
        refetch()
      }
    } finally {
      setSubmitting(prev => ({ ...prev, [category]: false }))
    }
  }

  const handleVote = async (itemId: string) => {
    try {
      const res = await fetch(`/api/retro/items/${itemId}/vote`, { method: 'POST' })
      if (res.ok) {
        refetch()
      }
    } catch {
      // silently fail
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent, category: RetroItem['category']) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleAddItem(category)
    }
  }

  if (retroLoading) {
    return (
      <div className="loading">
        <div className="loading-spinner" />
        Loading...
      </div>
    )
  }

  const hcName = hcData?.healthcheck?.Name || 'Health Check'
  const session = retroData?.session
  const items = retroData?.items || []

  const itemsByCategory = (category: RetroItem['category']) =>
    items
      .filter(item => item.category === category)
      .sort((a, b) => b.votes - a.votes)

  return (
    <div>
      <Link to={hcId ? `/healthcheck/${hcId}` : '/'} className="back-link">
        &#8592; Back to results
      </Link>

      <div className="page-header">
        <div className="page-header-top">
          <div>
            <h1>Retrospective</h1>
            <p style={{ color: 'var(--text-secondary)', fontSize: '14px', marginTop: '4px' }}>
              {hcName}
            </p>
          </div>
        </div>
      </div>

      {!session ? (
        <div className="glass-card" style={{ textAlign: 'center', padding: '48px 20px' }}>
          <div style={{ fontSize: '40px', marginBottom: '12px', opacity: 0.3 }}>
            {'\uD83D\uDD04'}
          </div>
          <p style={{ color: 'var(--text-secondary)', marginBottom: '20px', fontSize: '15px' }}>
            No retrospective session yet. Start one to collect feedback from the team.
          </p>
          <button
            className="btn btn-primary btn-lg"
            onClick={handleStartRetro}
            disabled={submitting.start}
          >
            Start Retrospective
          </button>
        </div>
      ) : (
        <div className="retro-columns">
          {CATEGORIES.map(({ key, label, color, dim }) => {
            const categoryItems = itemsByCategory(key)
            return (
              <div key={key} className="retro-column">
                <div
                  className="retro-column-header"
                  style={{ borderLeftColor: color }}
                >
                  <span className="retro-column-title">{label}</span>
                  <span
                    className="retro-column-count"
                    style={{ background: dim, color }}
                  >
                    {categoryItems.length}
                  </span>
                </div>

                <div className="retro-items">
                  {categoryItems.map(item => (
                    <div key={item.id} className="glass-card retro-item">
                      <p className="retro-item-text">{item.text}</p>
                      <div className="retro-item-footer">
                        <div className="retro-item-author">
                          <div
                            className="retro-item-avatar"
                            style={{ backgroundColor: getAvatarColor(item.author) }}
                          >
                            {getInitial(item.author)}
                          </div>
                          <span className="retro-item-author-name">{item.author}</span>
                        </div>
                        <button
                          className="retro-vote-btn"
                          onClick={() => handleVote(item.id)}
                          aria-label={`Upvote: ${item.votes} votes`}
                        >
                          <span>{'\uD83D\uDC4D'}</span>
                          <span className="retro-vote-count">{item.votes}</span>
                        </button>
                      </div>
                    </div>
                  ))}

                  {categoryItems.length === 0 && (
                    <div className="retro-empty">
                      No items yet
                    </div>
                  )}
                </div>

                <div className="retro-add-form">
                  <input
                    ref={(el) => { formRefs.current[key] = el }}
                    type="text"
                    className="form-input retro-add-input"
                    placeholder={`Add ${label.toLowerCase()}...`}
                    value={newItemText[key]}
                    onChange={e => setNewItemText(prev => ({ ...prev, [key]: e.target.value }))}
                    onKeyDown={e => handleKeyDown(e, key)}
                    disabled={submitting[key]}
                  />
                  <button
                    className="btn btn-secondary btn-sm"
                    onClick={() => handleAddItem(key)}
                    disabled={!newItemText[key]?.trim() || submitting[key]}
                  >
                    Add
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
