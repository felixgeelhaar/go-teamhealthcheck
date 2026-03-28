package storage

import (
	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
	"github.com/felixgeelhaar/go-teamhealthcheck/sdk"
)

// SDKStoreReader adapts *Store to sdk.StoreReader by converting domain types to SDK types.
type SDKStoreReader struct {
	store *Store
}

// NewSDKStoreReader wraps a Store as an sdk.StoreReader.
func NewSDKStoreReader(s *Store) *SDKStoreReader {
	return &SDKStoreReader{store: s}
}

func (r *SDKStoreReader) FindTeamByID(id string) (*sdk.Team, error) {
	t, err := r.store.FindTeamByID(id)
	if err != nil || t == nil {
		return nil, err
	}
	return convertTeam(t), nil
}

func (r *SDKStoreReader) FindAllTeams() ([]*sdk.Team, error) {
	teams, err := r.store.FindAllTeams()
	if err != nil {
		return nil, err
	}
	result := make([]*sdk.Team, len(teams))
	for i, t := range teams {
		result[i] = convertTeam(t)
	}
	return result, nil
}

func (r *SDKStoreReader) FindHealthCheckByID(id string) (*sdk.HealthCheck, error) {
	hc, err := r.store.FindHealthCheckByID(id)
	if err != nil || hc == nil {
		return nil, err
	}
	return convertHealthCheck(hc), nil
}

func (r *SDKStoreReader) FindAllHealthChecks(filter sdk.HealthCheckFilter) ([]*sdk.HealthCheck, error) {
	df := domain.HealthCheckFilter{Limit: filter.Limit}
	if filter.TeamID != nil {
		df.TeamID = filter.TeamID
	}
	if filter.Status != nil {
		s := domain.Status(*filter.Status)
		df.Status = &s
	}
	hcs, err := r.store.FindAllHealthChecks(df)
	if err != nil {
		return nil, err
	}
	result := make([]*sdk.HealthCheck, len(hcs))
	for i, hc := range hcs {
		result[i] = convertHealthCheck(hc)
	}
	return result, nil
}

func (r *SDKStoreReader) FindTemplateByID(id string) (*sdk.Template, error) {
	t, err := r.store.FindTemplateByID(id)
	if err != nil || t == nil {
		return nil, err
	}
	return convertTemplate(t), nil
}

func (r *SDKStoreReader) FindVotesByHealthCheck(id string) ([]*sdk.Vote, error) {
	votes, err := r.store.FindVotesByHealthCheck(id)
	if err != nil {
		return nil, err
	}
	result := make([]*sdk.Vote, len(votes))
	for i, v := range votes {
		result[i] = &sdk.Vote{
			ID:            v.ID,
			HealthCheckID: v.HealthCheckID,
			MetricName:    v.MetricName,
			Participant:   v.Participant,
			Color:         sdk.VoteColor(v.Color),
			Comment:       v.Comment,
			CreatedAt:     v.CreatedAt,
		}
	}
	return result, nil
}

func convertTeam(t *domain.Team) *sdk.Team {
	return &sdk.Team{
		ID:        t.ID,
		Name:      t.Name,
		Members:   t.Members,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func convertHealthCheck(hc *domain.HealthCheck) *sdk.HealthCheck {
	return &sdk.HealthCheck{
		ID:         hc.ID,
		TeamID:     hc.TeamID,
		TemplateID: hc.TemplateID,
		Name:       hc.Name,
		Anonymous:  hc.Anonymous,
		Status:     string(hc.Status),
		CreatedAt:  hc.CreatedAt,
		ClosedAt:   hc.ClosedAt,
	}
}

func convertTemplate(t *domain.Template) *sdk.Template {
	metrics := make([]sdk.TemplateMetric, len(t.Metrics))
	for i, m := range t.Metrics {
		metrics[i] = sdk.TemplateMetric{
			ID:              m.ID,
			TemplateID:      m.TemplateID,
			Name:            m.Name,
			DescriptionGood: m.DescriptionGood,
			DescriptionBad:  m.DescriptionBad,
			SortOrder:       m.SortOrder,
		}
	}
	return &sdk.Template{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		BuiltIn:     t.BuiltIn,
		Metrics:     metrics,
		CreatedAt:   t.CreatedAt,
	}
}
