package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"
	"github.com/felixgeelhaar/go-teamhealthcheck/internal/storage"
)

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func handleAPITeams(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := store.FindAllTeams()
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if teams == nil {
			teams = []*domain.Team{}
		}
		writeJSON(w, teams)
	}
}

func handleAPIHealthChecks(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := domain.HealthCheckFilter{Limit: 50}
		if teamID := r.URL.Query().Get("team_id"); teamID != "" {
			filter.TeamID = &teamID
		}
		if status := r.URL.Query().Get("status"); status != "" {
			s := domain.Status(status)
			filter.Status = &s
		}

		hcs, err := store.FindAllHealthChecks(filter)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if hcs == nil {
			hcs = []*domain.HealthCheck{}
		}
		writeJSON(w, hcs)
	}
}

func handleAPIHealthCheck(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		hc, err := store.FindHealthCheckByID(id)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if hc == nil {
			writeError(w, 404, "health check not found")
			return
		}

		tmpl, err := store.FindTemplateByID(hc.TemplateID)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		writeJSON(w, map[string]any{
			"healthcheck": hc,
			"template":    tmpl,
		})
	}
}

func handleAPIResults(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		hc, err := store.FindHealthCheckByID(id)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if hc == nil {
			writeError(w, 404, "health check not found")
			return
		}

		tmpl, err := store.FindTemplateByID(hc.TemplateID)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		votes, err := store.FindVotesByHealthCheck(id)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		results := domain.ComputeMetricResults(votes, tmpl.Metrics)

		var totalScore float64
		var totalVotes int
		participants := make(map[string]bool)
		for _, v := range votes {
			participants[v.Participant] = true
		}
		for _, res := range results {
			totalScore += res.Score * float64(res.TotalVotes)
			totalVotes += res.TotalVotes
		}
		var avgScore float64
		if totalVotes > 0 {
			avgScore = totalScore / float64(totalVotes)
		}

		writeJSON(w, map[string]any{
			"healthcheck":   hc,
			"results":       results,
			"average_score": avgScore,
			"participants":  len(participants),
			"total_votes":   totalVotes,
		})
	}
}
