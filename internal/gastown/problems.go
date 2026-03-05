package gastown

// Problem represents a detected issue with a Gas Town agent or beads infrastructure.
type Problem struct {
	Type     string       // "stalled", "backoff", "zombie", "doctor"
	Agent    AgentRuntime // the affected agent (zero value for doctor problems)
	Detail   string       // human-readable description
	Severity string       // "warn", "error"
	Category string       // doctor category (e.g. "Core System", "Git Integration")
	Fix      string       // suggested fix command, if any
}

// DetectProblems analyzes TownStatus and returns any detected problems.
// All heuristics are stateless — derived from a single status snapshot.
func DetectProblems(status *TownStatus) []Problem {
	if status == nil {
		return nil
	}

	var problems []Problem
	for _, a := range status.Agents {
		// Stalled: agent is running, has work, but idle (should be working)
		if a.HasWork && a.State == "idle" {
			problems = append(problems, Problem{
				Type:     "stalled",
				Agent:    a,
				Detail:   "Has work but idle — may need nudge",
				Severity: "warn",
			})
		}

		// Stuck: agent explicitly requesting help
		if a.State == "stuck" {
			problems = append(problems, Problem{
				Type:     "stuck",
				Agent:    a,
				Detail:   "Agent is stuck and requesting help",
				Severity: "error",
			})
		}

		// Backoff spiral: agent is in backoff state
		if a.State == "backoff" {
			problems = append(problems, Problem{
				Type:     "backoff",
				Agent:    a,
				Detail:   "In backoff state — may be stuck in retry loop",
				Severity: "warn",
			})
		}

		// Zombie: agent not running but has hooked work
		if !a.Running && a.HookBead != "" {
			problems = append(problems, Problem{
				Type:     "zombie",
				Agent:    a,
				Detail:   "Not running but has hooked work (" + a.HookBead + ")",
				Severity: "error",
			})
		}
	}

	return problems
}

// DoctorProblems converts bd doctor diagnostics into Problem entries.
// Only error and warning diagnostics are included (not passed checks).
func DoctorProblems(diagnostics []DoctorDiagnostic) []Problem {
	var problems []Problem
	for _, d := range diagnostics {
		if d.Status == "ok" {
			continue
		}
		sev := "warn"
		if d.Status == "error" {
			sev = "error"
		}
		fix := ""
		if len(d.Commands) > 0 {
			fix = d.Commands[0]
		}
		problems = append(problems, Problem{
			Type:     "doctor",
			Detail:   d.Name + ": " + d.Explanation,
			Severity: sev,
			Category: d.Category,
			Fix:      fix,
		})
	}
	return problems
}

// DoctorDiagnostic mirrors data.DoctorDiagnostic for use within gastown package.
type DoctorDiagnostic struct {
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Category    string   `json:"category"`
	Explanation string   `json:"explanation"`
	Commands    []string `json:"commands"`
}
