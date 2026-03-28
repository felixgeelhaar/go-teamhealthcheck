package mcpui

import (
	"fmt"
	"html"
	"strings"

	"github.com/felixgeelhaar/heartbeat/internal/domain"
)

// VotingFormHTML generates a self-contained HTML voting form for MCP Apps.
// Renders inside a sandboxed iframe in Claude Desktop.
func VotingFormHTML(healthCheckID string, metrics []domain.TemplateMetric) string {
	var metricCards strings.Builder
	for _, m := range metrics {
		metricCards.WriteString(fmt.Sprintf(`
		<div class="metric">
			<div class="metric-name">%s</div>
			<div class="descriptions">
				<span class="good">%s</span>
				<span class="bad">%s</span>
			</div>
			<div class="vote-buttons">
				<label class="vote green"><input type="radio" name="%s" value="green"><span></span></label>
				<label class="vote yellow"><input type="radio" name="%s" value="yellow"><span></span></label>
				<label class="vote red"><input type="radio" name="%s" value="red"><span></span></label>
			</div>
		</div>`,
			html.EscapeString(m.Name),
			html.EscapeString(m.DescriptionGood),
			html.EscapeString(m.DescriptionBad),
			html.EscapeString(m.Name),
			html.EscapeString(m.Name),
			html.EscapeString(m.Name),
		))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;padding:16px;color:#111827;background:#fff}
h2{font-size:18px;margin-bottom:16px}
.metric{border:1px solid #e5e7eb;border-radius:8px;padding:12px;margin-bottom:12px}
.metric-name{font-weight:700;font-size:15px;margin-bottom:6px}
.descriptions{display:flex;justify-content:space-between;font-size:11px;margin-bottom:10px;gap:8px}
.good{color:#166534;flex:1}.bad{color:#991b1b;flex:1;text-align:right}
.vote-buttons{display:flex;gap:12px;justify-content:center}
.vote{cursor:pointer;display:flex;align-items:center}
.vote input{display:none}
.vote span{width:32px;height:32px;border-radius:50%%;border:3px solid #d1d5db;transition:all .15s}
.vote.green span{border-color:#86efac}.vote.yellow span{border-color:#fde68a}.vote.red span{border-color:#fca5a5}
.vote input:checked~span{transform:scale(1.1)}
.vote.green input:checked~span{background:#22c55e;border-color:#22c55e}
.vote.yellow input:checked~span{background:#eab308;border-color:#eab308}
.vote.red input:checked~span{background:#ef4444;border-color:#ef4444}
.submit{display:block;width:100%%;padding:12px;background:#3b82f6;color:#fff;border:none;border-radius:8px;font-size:15px;font-weight:600;cursor:pointer;margin-top:16px}
.submit:hover{background:#2563eb}
.submit:disabled{background:#9ca3af;cursor:not-allowed}
</style></head><body>
<h2>Health Check Vote</h2>
%s
<button class="submit" onclick="submitVotes()">Submit Votes</button>
<script>
function submitVotes(){
  var votes=[];
  document.querySelectorAll('.metric').forEach(function(m){
    var name=m.querySelector('.metric-name').textContent;
    var checked=m.querySelector('input:checked');
    if(checked)votes.push({metric_name:name,color:checked.value});
  });
  if(votes.length===0){alert('Please vote on at least one metric');return}
  window.parent.postMessage({type:'mcp-app-result',data:{healthcheck_id:'%s',votes:votes}},'*');
  document.querySelector('.submit').disabled=true;
  document.querySelector('.submit').textContent='Votes submitted!';
}
</script></body></html>`, metricCards.String(), html.EscapeString(healthCheckID))
}
