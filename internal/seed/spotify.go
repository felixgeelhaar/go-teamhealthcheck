package seed

import "github.com/felixgeelhaar/go-teamhealthcheck/internal/domain"

// SpotifyTemplate returns the original Spotify Squad Health Check template.
// Based on Henrik Kniberg's model: https://labs.spotify.com/2014/09/16/squad-health-check-model/
func SpotifyTemplate() domain.Template {
	return domain.Template{
		Name:        "Spotify Squad Health Check",
		Description: "The original Squad Health Check model by Henrik Kniberg at Spotify",
		BuiltIn:     true,
		Metrics: []domain.TemplateMetric{
			{Name: "Easy to Release", DescriptionGood: "Releasing is simple, safe, painless & mostly automated", DescriptionBad: "Releasing is risky, painful, lots of manual work, and takes forever", SortOrder: 1},
			{Name: "Suitable Process", DescriptionGood: "Our way of working fits us perfectly", DescriptionBad: "Our way of working sucks", SortOrder: 2},
			{Name: "Tech Quality", DescriptionGood: "We're proud of the quality of our code! It is clean, easy to read, and has great test coverage", DescriptionBad: "Our code is a pile of dung, and technical debt is raging out of control", SortOrder: 3},
			{Name: "Value", DescriptionGood: "We deliver great stuff! We're proud of it and our stakeholders are really happy", DescriptionBad: "We deliver crap. We feel ashamed to deliver it. Our stakeholders hate us", SortOrder: 4},
			{Name: "Speed", DescriptionGood: "We get stuff done really quickly. No waiting, no delays", DescriptionBad: "We never seem to get done with anything. Stories keep getting stuck on dependencies", SortOrder: 5},
			{Name: "Mission", DescriptionGood: "We know exactly why we are here, and we are really excited about it", DescriptionBad: "We have no idea why we are here, there is no clear mission", SortOrder: 6},
			{Name: "Fun", DescriptionGood: "We love going to work, and have great fun working together", DescriptionBad: "Boooooring", SortOrder: 7},
			{Name: "Learning", DescriptionGood: "We're learning lots of interesting stuff all the time!", DescriptionBad: "We never have time to learn anything", SortOrder: 8},
			{Name: "Support", DescriptionGood: "We always get great support & help when we ask for it!", DescriptionBad: "We keep getting stuck because we can't get the support & help we ask for", SortOrder: 9},
			{Name: "Pawns or Players", DescriptionGood: "We are in control of our destiny! We decide what to build and how to build it", DescriptionBad: "We are just pawns in a game of chess, with no influence over what we build or how we build it", SortOrder: 10},
		},
	}
}
