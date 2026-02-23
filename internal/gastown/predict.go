package gastown

import (
	"fmt"
	"time"

	"github.com/matt-wright86/mardi-gras/internal/data"
)

// ConvoyPrediction holds ETA and progress rate for a convoy.
type ConvoyPrediction struct {
	ConvoyID       string
	ConvoyTitle    string
	Completed      int
	Total          int
	Remaining      int
	AvgTimePerTask time.Duration // based on velocity
	EstimatedETA   time.Duration // estimated time to completion
	ETALabel       string        // human-readable ETA
	Confidence     string        // "high", "medium", "low"
}

// PredictConvoys estimates completion times for active convoys.
// Uses velocity metrics (closed-per-day rate) to project ETAs.
func PredictConvoys(convoys []ConvoyDetail, velocity *VelocityMetrics) []ConvoyPrediction {
	if len(convoys) == 0 || velocity == nil {
		return nil
	}

	// Compute daily close rate from velocity
	// Use week rate if today rate is too small (more stable)
	dailyRate := float64(velocity.ClosedWeek) / 7.0
	if velocity.ClosedToday > 0 {
		// Weighted blend: today's rate is more recent but volatile
		todayRate := float64(velocity.ClosedToday)
		dailyRate = todayRate*0.3 + dailyRate*0.7
	}

	var predictions []ConvoyPrediction

	for _, c := range convoys {
		if c.Status == "closed" || c.Total == 0 {
			continue
		}

		remaining := c.Total - c.Completed
		if remaining <= 0 {
			continue
		}

		pred := ConvoyPrediction{
			ConvoyID:    c.ID,
			ConvoyTitle: c.Title,
			Completed:   c.Completed,
			Total:       c.Total,
			Remaining:   remaining,
		}

		if dailyRate > 0 {
			daysRemaining := float64(remaining) / dailyRate
			pred.EstimatedETA = time.Duration(daysRemaining * float64(24*time.Hour))
			pred.ETALabel = formatETA(pred.EstimatedETA)

			// Confidence based on data quality
			switch {
			case velocity.ClosedWeek >= 5 && velocity.WorkingAgents > 0:
				pred.Confidence = "high"
			case velocity.ClosedWeek >= 2:
				pred.Confidence = "medium"
			default:
				pred.Confidence = "low"
			}
		} else {
			pred.ETALabel = "unknown"
			pred.Confidence = "low"
		}

		predictions = append(predictions, pred)
	}

	return predictions
}

// PredictCostBudget estimates when the daily cost budget will be exhausted.
func PredictCostBudget(costs *CostsOutput, dailyBudget float64, now time.Time) string {
	if costs == nil || dailyBudget <= 0 || costs.Total.Cost <= 0 {
		return ""
	}

	// How much time has elapsed today
	y, m, d := now.Date()
	todayStart := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	elapsed := now.Sub(todayStart)
	if elapsed <= 0 {
		return ""
	}

	// Cost burn rate (per hour)
	hourlyRate := costs.Total.Cost / elapsed.Hours()
	if hourlyRate <= 0 {
		return ""
	}

	// Time until budget exhausted
	budgetRemaining := dailyBudget - costs.Total.Cost
	if budgetRemaining <= 0 {
		return "budget exceeded"
	}

	hoursLeft := budgetRemaining / hourlyRate
	if hoursLeft < 1 {
		return fmt.Sprintf("~%dm left at current rate", int(hoursLeft*60))
	}
	return fmt.Sprintf("~%.1fh left at current rate", hoursLeft)
}

// AgentProductivity estimates issues-per-session for active agents.
func AgentProductivity(issues []data.Issue, costs *CostsOutput) float64 {
	if costs == nil || costs.Sessions == 0 {
		return 0
	}
	var closedToday int
	todayStart := time.Now().Truncate(24 * time.Hour)
	for _, iss := range issues {
		if iss.Status == data.StatusClosed && iss.ClosedAt != nil {
			if iss.ClosedAt.After(todayStart) || iss.ClosedAt.Equal(todayStart) {
				closedToday++
			}
		}
	}
	return float64(closedToday) / float64(costs.Sessions)
}

// formatETA returns a human-readable duration string.
func formatETA(d time.Duration) string {
	hours := d.Hours()
	switch {
	case hours < 1:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case hours < 24:
		return fmt.Sprintf("%.1fh", hours)
	case hours < 168: // 7 days
		return fmt.Sprintf("%.1fd", hours/24)
	default:
		return fmt.Sprintf("%.0fw", hours/168)
	}
}
