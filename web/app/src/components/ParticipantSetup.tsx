import { useState } from 'react'
import { useParticipant } from '../hooks/useParticipant'

interface Props {
  onComplete: () => void
}

export function ParticipantSetup({ onComplete }: Props) {
  const { save } = useParticipant()
  const [inputValue, setInputValue] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = inputValue.trim()
    if (trimmed) {
      save(trimmed)
      onComplete()
    }
  }

  return (
    <div className="overlay">
      <div className="overlay-card glass-card">
        <h2>Welcome</h2>
        <p>Enter your name to get started with health checks.</p>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <input
              type="text"
              className="form-input"
              placeholder="Your name"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              autoFocus
              style={{ textAlign: 'center', fontSize: '16px' }}
            />
          </div>
          <button
            type="submit"
            className="btn btn-primary btn-lg"
            disabled={!inputValue.trim()}
            style={{ width: '100%' }}
          >
            Continue
          </button>
        </form>
      </div>
    </div>
  )
}
