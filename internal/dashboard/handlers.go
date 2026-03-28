package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

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

// --- GET handlers ---

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

		avgScore, totalVotes, participantNames := domain.ComputeOverallScore(results, votes)

		actions, _ := store.FindActionsByHealthCheck(id)

		// Strip names and comments if anonymous
		if hc.Anonymous {
			participantNames = []string{}
			for i := range results {
				results[i].Comments = []string{}
			}
		}

		writeJSON(w, map[string]any{
			"healthcheck":       hc,
			"results":           results,
			"average_score":     avgScore,
			"participants":      len(participantNames),
			"participant_names": participantNames,
			"total_votes":       totalVotes,
			"actions":           actions,
		})
	}
}

func handleAPITeamTrends(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID := r.PathValue("id")

		hcs, err := store.FindAllHealthChecks(domain.HealthCheckFilter{
			TeamID: &teamID,
			Limit:  20,
		})
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		if len(hcs) == 0 {
			writeJSON(w, map[string]any{
				"team_id":  teamID,
				"sessions": []any{},
				"trends":   []any{},
			})
			return
		}

		// Reverse to chronological order (oldest first)
		for i, j := 0, len(hcs)-1; i < j; i, j = i+1, j-1 {
			hcs[i], hcs[j] = hcs[j], hcs[i]
		}

		// Build per-session summaries and per-metric trends
		type sessionSummary struct {
			ID       string  `json:"id"`
			Name     string  `json:"name"`
			Date     string  `json:"date"`
			AvgScore float64 `json:"avg_score"`
			Voters   int     `json:"voters"`
			Status   string  `json:"status"`
		}

		metricScores := make(map[string][]domain.SessionScore)
		var sessions []sessionSummary

		for _, hc := range hcs {
			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil {
				continue
			}
			votes, err := store.FindVotesByHealthCheck(hc.ID)
			if err != nil {
				continue
			}

			results := domain.ComputeMetricResults(votes, tmpl.Metrics)

			var totalScore float64
			var totalVotes int
			voterSet := make(map[string]bool)
			for _, v := range votes {
				voterSet[v.Participant] = true
			}
			for _, res := range results {
				totalScore += res.Score * float64(res.TotalVotes)
				totalVotes += res.TotalVotes

				metricScores[res.MetricName] = append(metricScores[res.MetricName], domain.SessionScore{
					HealthCheckID:   hc.ID,
					HealthCheckName: hc.Name,
					Score:           res.Score,
					Date:            hc.CreatedAt.Format("2006-01-02"),
				})
			}

			var avgScore float64
			if totalVotes > 0 {
				avgScore = totalScore / float64(totalVotes)
			}

			sessions = append(sessions, sessionSummary{
				ID:       hc.ID,
				Name:     hc.Name,
				Date:     hc.CreatedAt.Format("2006-01-02"),
				AvgScore: avgScore,
				Voters:   len(voterSet),
				Status:   string(hc.Status),
			})
		}

		// Compute trends
		var trends []domain.MetricTrend
		for name, scores := range metricScores {
			var delta float64
			if len(scores) >= 2 && scores[0].Score > 0 {
				delta = scores[len(scores)-1].Score - scores[0].Score
			}
			trends = append(trends, domain.MetricTrend{
				MetricName: name,
				Sessions:   scores,
				Tendency:   domain.ComputeTendency(delta),
				Delta:      delta,
			})
		}

		writeJSON(w, map[string]any{
			"team_id":  teamID,
			"sessions": sessions,
			"trends":   trends,
		})
	}
}

func handleAPIAlerts(store *storage.Store) http.HandlerFunc {
	type alert struct {
		Metric    string  `json:"metric"`
		Severity  string  `json:"severity"` // "warning" or "critical"
		Message   string  `json:"message"`
		Score     float64 `json:"current_score"`
		Trend     string  `json:"trend"` // "declining", "stable", "improving"
		Delta     float64 `json:"delta"`
		Predicted float64 `json:"predicted_score"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		teamID := r.PathValue("id")

		hcs, err := store.FindAllHealthChecks(domain.HealthCheckFilter{
			TeamID: &teamID,
			Limit:  10,
		})
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		if len(hcs) < 2 {
			writeJSON(w, map[string]any{"alerts": []alert{}, "message": "Need at least 2 health checks for predictions"})
			return
		}

		// Reverse to chronological order
		for i, j := 0, len(hcs)-1; i < j; i, j = i+1, j-1 {
			hcs[i], hcs[j] = hcs[j], hcs[i]
		}

		// Collect per-metric scores across sessions
		metricScores := make(map[string][]float64)
		for _, hc := range hcs {
			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil {
				continue
			}
			votes, err := store.FindVotesByHealthCheck(hc.ID)
			if err != nil || len(votes) == 0 {
				continue
			}
			for _, res := range domain.ComputeMetricResults(votes, tmpl.Metrics) {
				if res.TotalVotes > 0 {
					metricScores[res.MetricName] = append(metricScores[res.MetricName], res.Score)
				}
			}
		}

		var alerts []alert
		for name, scores := range metricScores {
			if len(scores) < 2 {
				continue
			}

			latest := scores[len(scores)-1]
			delta := latest - scores[0]
			tendency := domain.ComputeTendency(delta)

			// Simple linear prediction: extend the trend one step
			avgDelta := delta / float64(len(scores)-1)
			predicted := latest + avgDelta
			if predicted < 1.0 {
				predicted = 1.0
			}
			if predicted > 3.0 {
				predicted = 3.0
			}

			var a *alert

			// Critical: already low AND declining
			if latest < domain.AlertThresholdCritical && tendency == domain.TendencyDeclining {
				a = &alert{
					Metric:    name,
					Severity:  "critical",
					Message:   fmt.Sprintf("%s is critically low (%.1f) and still declining. Immediate attention needed.", name, latest),
					Score:     latest,
					Trend:     string(tendency),
					Delta:     delta,
					Predicted: predicted,
				}
			} else if tendency == domain.TendencyDeclining && predicted < domain.AlertThresholdWarning {
				// Warning: declining and predicted to go below 2.0
				a = &alert{
					Metric:    name,
					Severity:  "warning",
					Message:   fmt.Sprintf("%s is declining (%.1f → predicted %.1f). May need attention next sprint.", name, latest, predicted),
					Score:     latest,
					Trend:     string(tendency),
					Delta:     delta,
					Predicted: predicted,
				}
			} else if latest < domain.AlertThresholdWarning {
				// Warning: currently low even if stable
				a = &alert{
					Metric:    name,
					Severity:  "warning",
					Message:   fmt.Sprintf("%s has been consistently low (%.1f). Consider dedicated improvement efforts.", name, latest),
					Score:     latest,
					Trend:     string(tendency),
					Delta:     delta,
					Predicted: predicted,
				}
			}

			if a != nil {
				alerts = append(alerts, *a)
			}
		}

		if alerts == nil {
			alerts = []alert{}
		}

		writeJSON(w, map[string]any{"alerts": alerts})
	}
}

func handleAPIDiscussion(store *storage.Store) http.HandlerFunc {
	type topic struct {
		Priority   int      `json:"priority"`
		Metric     string   `json:"metric"`
		Score      float64  `json:"score"`
		Reason     string   `json:"reason"`
		DataPoints []string `json:"data_points"`
		Questions  []string `json:"suggested_questions"`
	}

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

		var topics []topic
		priority := 1

		for _, res := range results {
			if res.TotalVotes == 0 {
				continue
			}

			var reasons []string
			var dataPoints []string
			var questions []string

			// High disagreement
			if res.GreenCount > 0 && res.RedCount > 0 {
				spread := float64(res.GreenCount-res.RedCount) / float64(res.TotalVotes)
				if spread < 0.5 && spread > -0.5 {
					reasons = append(reasons, "high disagreement")
					dataPoints = append(dataPoints, fmt.Sprintf("Split vote: %d green, %d yellow, %d red", res.GreenCount, res.YellowCount, res.RedCount))
					questions = append(questions, fmt.Sprintf("What different experiences lead to such varied opinions on %s?", res.MetricName))
				}
			}

			// Low score
			if res.Score < 2.0 {
				reasons = append(reasons, "low score")
				dataPoints = append(dataPoints, fmt.Sprintf("Score: %.1f/3.0", res.Score))
				questions = append(questions, fmt.Sprintf("What specific changes would improve %s the most?", res.MetricName))
			}

			if len(reasons) > 0 {
				reason := ""
				for i, r := range reasons {
					if i > 0 {
						reason += " + "
					}
					reason += r
				}
				topics = append(topics, topic{
					Priority:   priority,
					Metric:     res.MetricName,
					Score:      res.Score,
					Reason:     reason,
					DataPoints: dataPoints,
					Questions:  questions,
				})
				priority++
			}
		}

		if topics == nil {
			topics = []topic{}
		}

		writeJSON(w, map[string]any{
			"healthcheck_id": id,
			"topics":         topics,
		})
	}
}

func handleAPITemplates(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates, err := store.FindAllTemplates()
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if templates == nil {
			templates = []*domain.Template{}
		}
		writeJSON(w, templates)
	}
}

// --- POST handlers ---

func handleAPIVote(store *storage.Store) http.HandlerFunc {
	type voteRequest struct {
		Participant string `json:"participant"`
		MetricName  string `json:"metric_name"`
		Color       string `json:"color"`
		Comment     string `json:"comment"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var req voteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid request body")
			return
		}

		if req.Participant == "" || req.MetricName == "" || req.Color == "" {
			writeError(w, 400, "participant, metric_name, and color are required")
			return
		}

		// Validate health check exists and is open
		hc, err := store.FindHealthCheckByID(id)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if hc == nil {
			writeError(w, 404, "health check not found")
			return
		}
		if !hc.IsVotable() {
			writeError(w, 400, fmt.Sprintf("health check is %s, not accepting votes", hc.Status))
			return
		}

		// Get template for metric validation
		tmpl, err := store.FindTemplateByID(hc.TemplateID)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		// Use aggregate root to validate and create vote
		vote, err := hc.CastVote(req.MetricName, req.Participant, req.Color, req.Comment, tmpl.Metrics)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		vote.ID = uuid.NewString()
		vote.CreatedAt = time.Now()

		if err := store.UpsertVote(vote); err != nil {
			writeError(w, 500, err.Error())
			return
		}

		writeJSON(w, map[string]string{"status": "ok"})
	}
}

func handleAPICreateTemplate(store *storage.Store) http.HandlerFunc {
	type metricDef struct {
		Name            string `json:"name"`
		DescriptionGood string `json:"description_good"`
		DescriptionBad  string `json:"description_bad"`
	}
	type templateRequest struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Metrics     []metricDef `json:"metrics"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req templateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid request body")
			return
		}

		if req.Name == "" || len(req.Metrics) == 0 {
			writeError(w, 400, "name and at least one metric are required")
			return
		}

		tmplID := uuid.NewString()
		metrics := make([]domain.TemplateMetric, len(req.Metrics))
		for i, m := range req.Metrics {
			metrics[i] = domain.TemplateMetric{
				ID:              uuid.NewString(),
				TemplateID:      tmplID,
				Name:            m.Name,
				DescriptionGood: m.DescriptionGood,
				DescriptionBad:  m.DescriptionBad,
				SortOrder:       i + 1,
			}
		}

		tmpl := &domain.Template{
			ID:          tmplID,
			Name:        req.Name,
			Description: req.Description,
			BuiltIn:     false,
			Metrics:     metrics,
			CreatedAt:   time.Now(),
		}

		if err := store.CreateTemplate(tmpl); err != nil {
			writeError(w, 500, err.Error())
			return
		}

		writeJSON(w, tmpl)
	}
}

func handleAPICreateHealthCheck(store *storage.Store) http.HandlerFunc {
	type hcRequest struct {
		Name       string `json:"name"`
		TemplateID string `json:"template_id"`
		Anonymous  bool   `json:"anonymous"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		teamID := r.PathValue("id")

		var req hcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid request body")
			return
		}

		if req.Name == "" || req.TemplateID == "" {
			writeError(w, 400, "name and template_id are required")
			return
		}

		// Validate team
		team, err := store.FindTeamByID(teamID)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if team == nil {
			writeError(w, 404, "team not found")
			return
		}

		// Validate template
		tmpl, err := store.FindTemplateByID(req.TemplateID)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if tmpl == nil {
			writeError(w, 404, "template not found")
			return
		}

		hc := &domain.HealthCheck{
			ID:         uuid.NewString(),
			TeamID:     teamID,
			TemplateID: req.TemplateID,
			Name:       req.Name,
			Anonymous:  req.Anonymous,
			Status:     domain.StatusOpen,
			CreatedAt:  time.Now(),
		}

		if err := store.CreateHealthCheck(hc); err != nil {
			writeError(w, 500, err.Error())
			return
		}

		writeJSON(w, map[string]any{
			"healthcheck": hc,
			"template":    tmpl,
		})
	}
}

// --- CSV Export ---

func handleAPIExport(store *storage.Store) http.HandlerFunc {
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

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", strings.ReplaceAll(hc.Name, " ", "_")))

		fmt.Fprintf(w, "Metric,Green,Yellow,Red,Total Votes,Score,Description Good,Description Bad\n")
		for _, res := range results {
			fmt.Fprintf(w, "%s,%d,%d,%d,%d,%.2f,\"%s\",\"%s\"\n",
				csvEscape(res.MetricName),
				res.GreenCount, res.YellowCount, res.RedCount,
				res.TotalVotes, res.Score,
				csvEscape(res.DescriptionGood),
				csvEscape(res.DescriptionBad),
			)
		}
	}
}

func csvEscape(s string) string {
	return strings.ReplaceAll(s, "\"", "\"\"")
}

// --- Cross-Team Comparison ---

func handleAPICompare(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teams, err := store.FindAllTeams()
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}

		type teamSummary struct {
			TeamID   string  `json:"team_id"`
			TeamName string  `json:"team_name"`
			HCName   string  `json:"healthcheck_name"`
			HCID     string  `json:"healthcheck_id"`
			AvgScore float64 `json:"avg_score"`
			Voters   int     `json:"voters"`
			Date     string  `json:"date"`
			Status   string  `json:"status"`
		}

		var summaries []teamSummary

		for _, team := range teams {
			hcs, err := store.FindAllHealthChecks(domain.HealthCheckFilter{
				TeamID: &team.ID,
				Limit:  1,
			})
			if err != nil || len(hcs) == 0 {
				continue
			}

			hc := hcs[0]
			tmpl, err := store.FindTemplateByID(hc.TemplateID)
			if err != nil {
				continue
			}
			votes, err := store.FindVotesByHealthCheck(hc.ID)
			if err != nil || len(votes) == 0 {
				continue
			}

			results := domain.ComputeMetricResults(votes, tmpl.Metrics)
			var totalScore float64
			var totalVotes int
			voterSet := make(map[string]bool)
			for _, v := range votes {
				voterSet[v.Participant] = true
			}
			for _, res := range results {
				totalScore += res.Score * float64(res.TotalVotes)
				totalVotes += res.TotalVotes
			}
			var avgScore float64
			if totalVotes > 0 {
				avgScore = totalScore / float64(totalVotes)
			}

			summaries = append(summaries, teamSummary{
				TeamID:   team.ID,
				TeamName: team.Name,
				HCName:   hc.Name,
				HCID:     hc.ID,
				AvgScore: avgScore,
				Voters:   len(voterSet),
				Date:     hc.CreatedAt.Format("2006-01-02"),
				Status:   string(hc.Status),
			})
		}

		if summaries == nil {
			summaries = []teamSummary{}
		}

		writeJSON(w, map[string]any{"teams": summaries})
	}
}

// --- Action Items ---

func handleAPICreateAction(store *storage.Store) http.HandlerFunc {
	type actionRequest struct {
		MetricName  string `json:"metric_name"`
		Description string `json:"description"`
		Assignee    string `json:"assignee"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var req actionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "invalid request body")
			return
		}
		if req.Description == "" {
			writeError(w, 400, "description is required")
			return
		}

		hc, err := store.FindHealthCheckByID(id)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		if hc == nil {
			writeError(w, 404, "health check not found")
			return
		}

		action := &domain.Action{
			ID:            uuid.NewString(),
			HealthCheckID: id,
			MetricName:    req.MetricName,
			Description:   req.Description,
			Assignee:      req.Assignee,
			Completed:     false,
			CreatedAt:     time.Now(),
		}

		if err := store.CreateAction(action); err != nil {
			writeError(w, 500, err.Error())
			return
		}

		writeJSON(w, action)
	}
}

func handleAPIGenerateActions(store *storage.Store) http.HandlerFunc {
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

		// Use domain service to generate suggested actions
		suggested := domain.GenerateSuggestedActions(results, id)
		var generated []*domain.Action
		for _, action := range suggested {
			action.ID = uuid.NewString()
			action.CreatedAt = time.Now()
			if err := store.CreateAction(action); err == nil {
				generated = append(generated, action)
			}
		}

		if generated == nil {
			generated = []*domain.Action{}
		}

		writeJSON(w, map[string]any{
			"generated": len(generated),
			"actions":   generated,
		})
	}
}

func handleAPICompleteAction(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := store.CompleteAction(id); err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, map[string]string{"status": "completed", "action_id": id})
	}
}
