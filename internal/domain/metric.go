package domain

// MetricResult is a computed value object that aggregates votes for a single metric.
type MetricResult struct {
	MetricName      string
	DescriptionGood string
	DescriptionBad  string
	GreenCount      int
	YellowCount     int
	RedCount        int
	TotalVotes      int
	Score           float64
	Comments        []string
}

// ComputeMetricResults aggregates raw votes into per-metric results, enriched with template descriptions.
func ComputeMetricResults(votes []*Vote, templateMetrics []TemplateMetric) []MetricResult {
	// Index template metrics by name for O(1) lookup
	tmByName := make(map[string]TemplateMetric, len(templateMetrics))
	order := make(map[string]int, len(templateMetrics))
	for _, tm := range templateMetrics {
		tmByName[tm.Name] = tm
		order[tm.Name] = tm.SortOrder
	}

	type accumulator struct {
		green, yellow, red int
		total              float64
		comments           []string
	}

	accum := make(map[string]*accumulator)
	for _, v := range votes {
		a, ok := accum[v.MetricName]
		if !ok {
			a = &accumulator{}
			accum[v.MetricName] = a
		}
		switch v.Color {
		case VoteGreen:
			a.green++
		case VoteYellow:
			a.yellow++
		case VoteRed:
			a.red++
		}
		a.total += v.Color.Score()
		if v.Comment != "" {
			a.comments = append(a.comments, v.Comment)
		}
	}

	// Build results for all template metrics (even those with no votes)
	results := make([]MetricResult, 0, len(templateMetrics))
	for _, tm := range templateMetrics {
		r := MetricResult{
			MetricName:      tm.Name,
			DescriptionGood: tm.DescriptionGood,
			DescriptionBad:  tm.DescriptionBad,
		}
		if a, ok := accum[tm.Name]; ok {
			r.GreenCount = a.green
			r.YellowCount = a.yellow
			r.RedCount = a.red
			r.TotalVotes = a.green + a.yellow + a.red
			if r.TotalVotes > 0 {
				r.Score = a.total / float64(r.TotalVotes)
			}
			r.Comments = a.comments
		}
		if r.Comments == nil {
			r.Comments = []string{}
		}
		results = append(results, r)
	}

	return results
}

// Tendency is a value object describing the direction of change for a metric over time.
type Tendency string

const (
	TendencyImproving Tendency = "improving"
	TendencyStable    Tendency = "stable"
	TendencyDeclining Tendency = "declining"
)

// SessionScore captures a metric's score at a point in time.
type SessionScore struct {
	HealthCheckID   string
	HealthCheckName string
	Score           float64
	Date            string
}

// MetricTrend tracks how a single metric has changed across sessions.
type MetricTrend struct {
	MetricName string
	Sessions   []SessionScore
	Tendency   Tendency
	Delta      float64
}

// ComputeTendency derives a tendency from a delta value.
func ComputeTendency(delta float64) Tendency {
	switch {
	case delta > 0.2:
		return TendencyImproving
	case delta < -0.2:
		return TendencyDeclining
	default:
		return TendencyStable
	}
}
