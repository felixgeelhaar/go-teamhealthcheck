package seed

import "github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"

// TuckmanTemplate returns a team maturity assessment based on Tuckman's stages of group development.
func TuckmanTemplate() domain.Template {
	return domain.Template{
		Name:        "Team Maturity (Tuckman)",
		Description: "Assess team maturity across Tuckman's stages: Forming, Storming, Norming, Performing. Measures how well the team collaborates, resolves conflict, and self-organizes.",
		BuiltIn:     true,
		Metrics: []domain.TemplateMetric{
			{Name: "Trust & Safety", DescriptionGood: "Team members trust each other completely and feel safe to be vulnerable", DescriptionBad: "People guard information, avoid conflict, and don't trust each other's intentions", SortOrder: 1},
			{Name: "Healthy Conflict", DescriptionGood: "We engage in productive debate around ideas without making it personal", DescriptionBad: "Either we avoid all conflict or it gets personal and destructive", SortOrder: 2},
			{Name: "Commitment", DescriptionGood: "Everyone is committed to team decisions, even if they initially disagreed", DescriptionBad: "People say yes in meetings but don't follow through, or passively resist decisions", SortOrder: 3},
			{Name: "Accountability", DescriptionGood: "Team members hold each other accountable without needing the manager to step in", DescriptionBad: "Nobody calls out poor performance or missed commitments — it's left to the manager", SortOrder: 4},
			{Name: "Results Focus", DescriptionGood: "The team prioritizes collective results over individual goals and recognition", DescriptionBad: "People optimize for their own status, career, or department over team outcomes", SortOrder: 5},
			{Name: "Self-Organization", DescriptionGood: "The team organizes its own work, makes decisions autonomously, and adapts quickly", DescriptionBad: "Everything needs approval from above. The team waits to be told what to do", SortOrder: 6},
			{Name: "Role Clarity", DescriptionGood: "Everyone knows their role, responsibilities, and how they contribute to the team's goals", DescriptionBad: "Roles are unclear, work falls through the cracks, and people step on each other's toes", SortOrder: 7},
			{Name: "Continuous Improvement", DescriptionGood: "We regularly reflect on how we work and make concrete improvements", DescriptionBad: "We keep doing things the same way even when they're clearly not working", SortOrder: 8},
		},
	}
}

// PsychologicalSafetyTemplate returns an assessment based on Amy Edmondson's research on psychological safety.
func PsychologicalSafetyTemplate() domain.Template {
	return domain.Template{
		Name:        "Psychological Safety (Edmondson)",
		Description: "Based on Amy Edmondson's research at Harvard. Measures whether team members feel safe to take interpersonal risks — speak up, ask questions, admit mistakes, and propose new ideas.",
		BuiltIn:     true,
		Metrics: []domain.TemplateMetric{
			{Name: "Speaking Up", DescriptionGood: "Everyone feels comfortable voicing opinions, concerns, and disagreements openly", DescriptionBad: "People stay quiet in meetings to avoid being judged or shut down", SortOrder: 1},
			{Name: "Asking for Help", DescriptionGood: "It's normal and encouraged to ask for help when you're stuck or unsure", DescriptionBad: "Asking for help is seen as a sign of weakness or incompetence", SortOrder: 2},
			{Name: "Admitting Mistakes", DescriptionGood: "Mistakes are treated as learning opportunities and shared openly", DescriptionBad: "People hide mistakes or blame others to protect themselves", SortOrder: 3},
			{Name: "Risk Taking", DescriptionGood: "It's safe to experiment, try new approaches, and even fail", DescriptionBad: "People stick to the safe path because failure is punished or judged", SortOrder: 4},
			{Name: "Diversity of Thought", DescriptionGood: "Different perspectives, backgrounds, and ideas are actively sought and valued", DescriptionBad: "There's pressure to conform. Dissenting views are dismissed or ignored", SortOrder: 5},
			{Name: "Feedback Culture", DescriptionGood: "Giving and receiving feedback is a regular, constructive part of how we work", DescriptionBad: "Feedback is rare, vague, or only comes during formal reviews", SortOrder: 6},
			{Name: "Inclusion", DescriptionGood: "Everyone has equal opportunity to contribute and is heard in discussions", DescriptionBad: "Some voices dominate while others are overlooked or excluded", SortOrder: 7},
		},
	}
}

// DORATemplate returns an assessment based on DORA (DevOps Research and Assessment) metrics.
func DORATemplate() domain.Template {
	return domain.Template{
		Name:        "DORA Metrics (DevOps)",
		Description: "Based on the DORA research program's four key metrics for software delivery performance. Adapted as a team perception health check rather than quantitative measurement.",
		BuiltIn:     true,
		Metrics: []domain.TemplateMetric{
			{Name: "Deployment Frequency", DescriptionGood: "We deploy to production frequently and confidently — multiple times per day or week", DescriptionBad: "Deployments are rare, scary events that require weeks of planning and coordination", SortOrder: 1},
			{Name: "Lead Time for Changes", DescriptionGood: "Changes go from commit to production quickly — within hours or a day", DescriptionBad: "It takes weeks or months for a committed change to reach production", SortOrder: 2},
			{Name: "Change Failure Rate", DescriptionGood: "Our changes rarely cause incidents. When they do, we learn and prevent recurrence", DescriptionBad: "A large percentage of our changes cause failures, outages, or require hotfixes", SortOrder: 3},
			{Name: "Time to Restore", DescriptionGood: "When something breaks in production, we detect and fix it within minutes to hours", DescriptionBad: "Production incidents take days or weeks to resolve. Detection is slow", SortOrder: 4},
			{Name: "Test Confidence", DescriptionGood: "Our test suite gives us high confidence that changes work correctly", DescriptionBad: "Tests are flaky, incomplete, or don't exist. We ship and pray", SortOrder: 5},
			{Name: "Monitoring & Observability", DescriptionGood: "We have great visibility into system health, can quickly debug issues, and get alerted proactively", DescriptionBad: "We find out about problems from users. Debugging is guesswork", SortOrder: 6},
			{Name: "Infrastructure as Code", DescriptionGood: "Infrastructure is versioned, reproducible, and managed through code", DescriptionBad: "Servers are manually configured snowflakes. Nobody remembers how they were set up", SortOrder: 7},
			{Name: "Developer Experience", DescriptionGood: "Local dev setup is fast, CI/CD is reliable, tooling helps rather than hinders", DescriptionBad: "Dev environments are brittle, builds are slow, tooling is frustrating", SortOrder: 8},
		},
	}
}

// AllBuiltInTemplates returns all built-in assessment templates.
func AllBuiltInTemplates() []domain.Template {
	return []domain.Template{
		SpotifyTemplate(),
		TuckmanTemplate(),
		PsychologicalSafetyTemplate(),
		DORATemplate(),
	}
}
