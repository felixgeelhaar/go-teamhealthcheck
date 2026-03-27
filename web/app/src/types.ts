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
  Anonymous: boolean
  Status: 'open' | 'closed' | 'archived'
  CreatedAt: string
  ClosedAt: string | null
}

export interface TemplateMetric {
  ID: string
  TemplateID: string
  Name: string
  DescriptionGood: string
  DescriptionBad: string
  SortOrder: number
}

export interface Template {
  ID: string
  Name: string
  Description: string
  BuiltIn: boolean
  Metrics: TemplateMetric[]
  CreatedAt: string
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
  participant_names: string[]
  total_votes: number
}

export interface HealthCheckDetail {
  healthcheck: HealthCheck
  template: Template
}

export interface VotePayload {
  participant: string
  metric_name: string
  color: 'green' | 'yellow' | 'red'
  comment?: string
}

export interface CreateTemplateMetric {
  name: string
  description_good: string
  description_bad: string
}

export interface CreateTemplatePayload {
  name: string
  description: string
  metrics: CreateTemplateMetric[]
}

export interface CreateHealthCheckPayload {
  name: string
  template_id: string
  anonymous?: boolean
}

export interface DiscussionTopic {
  priority: number
  metric: string
  score: number
  reason: string
  data_points: string[]
  suggested_questions: string[]
}

export interface DiscussionGuideResponse {
  healthcheck_id: string
  topics: DiscussionTopic[]
}

export interface WSEvent {
  type: 'vote_submitted' | 'healthcheck_created' | 'healthcheck_status_changed' | 'healthcheck_deleted'
  healthcheck_id: string
  team_id?: string
  participant?: string
  metric_name?: string
}
