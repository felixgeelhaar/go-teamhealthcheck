export interface Team {
  ID: string
  Name: string
  Members: string[]
  CreatedAt: string
  UpdatedAt: string
}

export interface HealthCheck {
  ID: string
  TeamID: string
  TemplateID: string
  Name: string
  Status: 'open' | 'closed' | 'archived'
  CreatedAt: string
  ClosedAt: string | null
}

export interface MetricResult {
  MetricName: string
  DescriptionGood: string
  DescriptionBad: string
  GreenCount: number
  YellowCount: number
  RedCount: number
  TotalVotes: number
  Score: number
  Comments: string[]
}

export interface HealthCheckResults {
  healthcheck: HealthCheck
  results: MetricResult[]
  average_score: number
  participants: number
  total_votes: number
}

export interface WSEvent {
  type: 'vote_submitted' | 'healthcheck_created' | 'healthcheck_status_changed' | 'healthcheck_deleted'
  healthcheck_id: string
  team_id?: string
  participant?: string
  metric_name?: string
}
