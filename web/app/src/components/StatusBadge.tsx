interface Props {
  status: string
}

export function StatusBadge({ status }: Props) {
  const className = `status-badge status-badge-${status}`

  return (
    <span className={className}>
      <span className="badge-dot" />
      {status}
    </span>
  )
}
