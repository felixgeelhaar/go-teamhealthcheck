import type { MetricResult } from '../types'

interface Props {
  results: MetricResult[]
}

export function RadarChart({ results }: Props) {
  if (results.length < 3) return null

  const size = 400
  const center = size / 2
  const maxRadius = 150
  const labelOffset = 30
  const numAxes = results.length
  const rings = [1, 2, 3]

  const angleStep = (2 * Math.PI) / numAxes
  const startAngle = -Math.PI / 2

  function polarToCartesian(angle: number, radius: number): { x: number; y: number } {
    return {
      x: center + radius * Math.cos(angle),
      y: center + radius * Math.sin(angle),
    }
  }

  function getScoreColor(score: number): string {
    if (score >= 2.5) return 'var(--green)'
    if (score >= 1.5) return 'var(--yellow)'
    return 'var(--red)'
  }

  const ringPaths = rings.map((ringValue) => {
    const r = (ringValue / 3) * maxRadius
    const points = Array.from({ length: numAxes }, (_, i) => {
      const angle = startAngle + i * angleStep
      return polarToCartesian(angle, r)
    })
    const d = points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ') + ' Z'
    return <path key={ringValue} d={d} fill="none" stroke="rgba(255,255,255,0.1)" strokeWidth="1" />
  })

  const axisLines = Array.from({ length: numAxes }, (_, i) => {
    const angle = startAngle + i * angleStep
    const end = polarToCartesian(angle, maxRadius)
    return (
      <line
        key={i}
        x1={center}
        y1={center}
        x2={end.x}
        y2={end.y}
        stroke="rgba(255,255,255,0.1)"
        strokeWidth="1"
      />
    )
  })

  const dataPoints = results.map((result, i) => {
    const angle = startAngle + i * angleStep
    const r = (result.Score / 3) * maxRadius
    return polarToCartesian(angle, r)
  })

  const polygonPath =
    dataPoints.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ') + ' Z'

  const labels = results.map((result, i) => {
    const angle = startAngle + i * angleStep
    const pos = polarToCartesian(angle, maxRadius + labelOffset)

    let textAnchor: 'start' | 'middle' | 'end' = 'middle'
    const cosAngle = Math.cos(angle)
    if (cosAngle > 0.1) textAnchor = 'start'
    else if (cosAngle < -0.1) textAnchor = 'end'

    const name = result.MetricName
    const displayName = name.length > 16 ? name.substring(0, 14) + '...' : name

    return (
      <text
        key={result.MetricName}
        x={pos.x}
        y={pos.y}
        textAnchor={textAnchor}
        dominantBaseline="central"
        fill="var(--text-secondary)"
        fontSize="11"
        fontFamily="var(--font-sans)"
      >
        {displayName}
      </text>
    )
  })

  const dots = results.map((result, i) => {
    const point = dataPoints[i]
    return (
      <circle
        key={result.MetricName}
        cx={point.x}
        cy={point.y}
        r="4"
        fill={getScoreColor(result.Score)}
        stroke="var(--bg-primary)"
        strokeWidth="1.5"
      />
    )
  })

  const ringLabels = rings.map((ringValue) => {
    const y = center - (ringValue / 3) * maxRadius
    return (
      <text
        key={ringValue}
        x={center + 4}
        y={y - 4}
        fill="var(--text-tertiary)"
        fontSize="9"
        fontFamily="var(--font-mono)"
      >
        {ringValue.toFixed(1)}
      </text>
    )
  })

  return (
    <div className="glass-card" style={{ marginBottom: '24px' }}>
      <div className="section-title">Overview</div>
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <svg
          viewBox={`0 0 ${size} ${size}`}
          width="100%"
          height="auto"
          style={{ maxWidth: '400px' }}
        >
          {ringPaths}
          {axisLines}
          <polygon
            points={dataPoints.map((p) => `${p.x},${p.y}`).join(' ')}
            fill="rgba(59,130,246,0.3)"
            stroke="rgba(59,130,246,0.8)"
            strokeWidth="2"
          />
          {/* Hidden path for the polygon outline to maintain stacking */}
          <path d={polygonPath} fill="none" stroke="transparent" />
          {dots}
          {labels}
          {ringLabels}
        </svg>
      </div>
    </div>
  )
}
