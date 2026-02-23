package gastown

import (
	"sort"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

// AgentScorecard aggregates quality metrics for a single agent.
type AgentScorecard struct {
	Name          string
	IssuesClosed  int
	AvgQuality    float32 // average QualityScore across scored closed issues
	TotalScored   int     // how many closed issues had a quality score
	Crystallizing int     // issues where Crystallizes == true
	Ephemeral     int     // issues where Crystallizes == false
}

// ComputeScorecards derives per-agent quality scorecards from closed issues.
// It maps issues to agents via Assignee or Creator.Name fields.
func ComputeScorecards(issues []data.Issue) []AgentScorecard {
	agents := make(map[string]*AgentScorecard)

	for _, iss := range issues {
		if iss.Status != data.StatusClosed {
			continue
		}

		// Determine the agent name: prefer Assignee, fall back to Creator
		name := iss.Assignee
		if name == "" && iss.Creator != nil {
			name = iss.Creator.Name
		}
		if name == "" {
			continue
		}

		sc, ok := agents[name]
		if !ok {
			sc = &AgentScorecard{Name: name}
			agents[name] = sc
		}

		sc.IssuesClosed++

		if iss.QualityScore != nil {
			sc.TotalScored++
			sc.AvgQuality += *iss.QualityScore
		}

		if iss.Crystallizes != nil {
			if *iss.Crystallizes {
				sc.Crystallizing++
			} else {
				sc.Ephemeral++
			}
		}
	}

	// Finalize averages
	result := make([]AgentScorecard, 0, len(agents))
	for _, sc := range agents {
		if sc.TotalScored > 0 {
			sc.AvgQuality /= float32(sc.TotalScored)
		}
		result = append(result, *sc)
	}

	// Sort by average quality descending, then by issues closed
	sort.Slice(result, func(i, j int) bool {
		if result[i].AvgQuality != result[j].AvgQuality {
			return result[i].AvgQuality > result[j].AvgQuality
		}
		return result[i].IssuesClosed > result[j].IssuesClosed
	})

	return result
}
