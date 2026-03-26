import { useState } from 'react'
import { VoteButtons } from './VoteButtons'
import type { TemplateMetric } from '../types'

type VoteColor = 'green' | 'yellow' | 'red'

interface Props {
  metric: TemplateMetric
  selectedColor: VoteColor | null
  onColorSelect: (color: VoteColor) => void
  comment: string
  onCommentChange: (comment: string) => void
}

export function MetricCard({ metric, selectedColor, onColorSelect, comment, onCommentChange }: Props) {
  const [showComment, setShowComment] = useState(false)

  return (
    <div className="glass-card metric-card">
      <div className="metric-card-title">{metric.Name}</div>

      <div className="metric-card-descriptions">
        <div className="metric-desc metric-desc-good">
          <div className="metric-desc-label">What good looks like</div>
          {metric.DescriptionGood}
        </div>
        <div className="metric-desc metric-desc-bad">
          <div className="metric-desc-label">What bad looks like</div>
          {metric.DescriptionBad}
        </div>
      </div>

      <VoteButtons selected={selectedColor} onSelect={onColorSelect} />

      <div
        className="comment-toggle"
        onClick={() => setShowComment(!showComment)}
      >
        {showComment ? '- Hide comment' : '+ Add comment (optional)'}
      </div>

      {showComment && (
        <div style={{ marginTop: '8px' }}>
          <textarea
            className="form-textarea"
            placeholder="Share your thoughts..."
            value={comment}
            onChange={(e) => onCommentChange(e.target.value)}
            rows={2}
          />
        </div>
      )}
    </div>
  )
}
