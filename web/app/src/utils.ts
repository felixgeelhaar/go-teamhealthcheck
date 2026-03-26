export function scoreColorClass(score: number): string {
  if (score >= 2.5) return 'score-green'
  if (score >= 1.5) return 'score-yellow'
  return 'score-red'
}

export function getAvatarColor(name: string): string {
  const colors = ['#3b82f6', '#a855f7', '#ec4899', '#f97316', '#14b8a6', '#6366f1']
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

export function getInitial(name: string): string {
  return name.charAt(0).toUpperCase()
}

export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
