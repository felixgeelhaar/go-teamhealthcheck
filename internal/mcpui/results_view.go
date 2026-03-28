package mcpui

import (
	"fmt"
	"html"
	"strings"

	"github.com/felixgeelhaar/heartbeat/internal/domain"
)

// ResultsViewHTML generates a self-contained HTML results heatmap for MCP Apps.
func ResultsViewHTML(results []domain.MetricResult, avgScore float64, participants int) string {
	var rows strings.Builder
	for _, r := range results {
		total := r.TotalVotes
		if total == 0 {
			total = 1
		}
		greenPct := float64(r.GreenCount) / float64(total) * 100
		yellowPct := float64(r.YellowCount) / float64(total) * 100
		redPct := float64(r.RedCount) / float64(total) * 100

		scoreColor := "#9ca3af"
		if r.TotalVotes > 0 {
			if r.Score >= 2.5 {
				scoreColor = "#22c55e"
			} else if r.Score >= 1.5 {
				scoreColor = "#eab308"
			} else {
				scoreColor = "#ef4444"
			}
		}

		scoreText := "-"
		if r.TotalVotes > 0 {
			scoreText = fmt.Sprintf("%.1f", r.Score)
		}

		rows.WriteString(fmt.Sprintf(`
		<div class="row">
			<div class="name">%s<span class="votes">%d votes</span></div>
			<div class="bar">
				<div class="g" style="width:%.1f%%"></div>
				<div class="y" style="width:%.1f%%"></div>
				<div class="r" style="width:%.1f%%"></div>
			</div>
			<div class="score" style="color:%s">%s</div>
		</div>`,
			html.EscapeString(r.MetricName), r.TotalVotes,
			greenPct, yellowPct, redPct,
			scoreColor, scoreText,
		))
	}

	avgColor := "#9ca3af"
	if avgScore >= 2.5 {
		avgColor = "#22c55e"
	} else if avgScore >= 1.5 {
		avgColor = "#eab308"
	} else if avgScore > 0 {
		avgColor = "#ef4444"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;padding:16px;color:#111827;background:#fff}
.header{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px}
h2{font-size:18px}
.avg{font-size:28px;font-weight:800}
.meta{font-size:13px;color:#6b7280;margin-bottom:16px}
.row{display:grid;grid-template-columns:160px 1fr 50px;gap:10px;align-items:center;padding:6px 0;border-bottom:1px solid #f3f4f6}
.name{font-size:13px;font-weight:600}
.votes{display:block;font-size:11px;color:#9ca3af;font-weight:400}
.bar{display:flex;height:20px;border-radius:3px;overflow:hidden;background:#f3f4f6}
.g{background:#22c55e;transition:width .3s}.y{background:#eab308;transition:width .3s}.r{background:#ef4444;transition:width .3s}
.score{font-weight:700;font-size:15px;text-align:center}
</style></head><body>
<div class="header">
	<h2>Results</h2>
	<div class="avg" style="color:%s">%.1f</div>
</div>
<div class="meta">%d participant%s</div>
%s
</body></html>`,
		avgColor, avgScore,
		participants, pluralS(participants),
		rows.String(),
	)
}

func pluralS(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}
