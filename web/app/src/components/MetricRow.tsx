import { useState } from 'react'
import type { MetricResult } from '../types'
import { scoreColorClass } from '../utils'

interface Props {
  metric: MetricResult
}

export function MetricRow({ metric }: Props) {
  const [expanded, setExpanded] = useState(false)

  const total = metric.TotalVotes || 1
  const greenPct = (metric.GreenCount / total) * 100
  const yellowPct = (metric.YellowCount / total) * 100
  const redPct = (metric.RedCount / total) * 100

  const colorClass = scoreColorClass(metric.Score)

  return (
    <div className="metric-row">
      <div className="metric-row-header" onClick={() => setExpanded(!expanded)}>
        <div>
          <div className="metric-name">
            {metric.MetricName}
            <span className={`metric-expand-icon ${expanded ? 'expanded' : ''}`}>
              &#9654;
            </span>
          </div>
          <div className="metric-votes-badges">
            <span className="vote-badge">
              <span className="vote-badge-dot" style={{ background: 'var(--green)' }} />
              {metric.GreenCount}
            </span>
            <span className="vote-badge">
              <span className="vote-badge-dot" style={{ background: 'var(--yellow)' }} />
              {metric.YellowCount}
            </span>
            <span className="vote-badge">
              <span className="vote-badge-dot" style={{ background: 'var(--red)' }} />
              {metric.RedCount}
            </span>
          </div>
        </div>

        <div className="stacked-bar">
          {metric.GreenCount > 0 && (
            <div className="stacked-bar-segment stacked-bar-green" style={{ width: `${greenPct}%` }} />
          )}
          {metric.YellowCount > 0 && (
            <div className="stacked-bar-segment stacked-bar-yellow" style={{ width: `${yellowPct}%` }} />
          )}
          {metric.RedCount > 0 && (
            <div className="stacked-bar-segment stacked-bar-red" style={{ width: `${redPct}%` }} />
          )}
        </div>

        <div className={`metric-score ${colorClass}`}>
          {metric.TotalVotes > 0 ? metric.Score.toFixed(1) : '-'}
        </div>
      </div>

      {expanded && (
        <div className="metric-expanded">
          <div className="metric-desc metric-desc-good">
            <div className="metric-desc-label">What good looks like</div>
            {metric.DescriptionGood}
          </div>
          <div className="metric-desc metric-desc-bad">
            <div className="metric-desc-label">What bad looks like</div>
            {metric.DescriptionBad}
          </div>
        </div>
      )}
    </div>
  )
}
