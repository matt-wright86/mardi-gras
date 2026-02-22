package data

import (
	"strings"
)

// FilterIssues returns a new slice of issues that match the search query.
// It supports explicit tokens (type:bug, p1, priority:high) and free-text search on ID and Title.
// All tokens in the query must match (AND logic) for an issue to be included.
func FilterIssues(issues []Issue, query string) []Issue {
	query = strings.TrimSpace(query)
	if query == "" {
		return issues
	}

	rawTokens := strings.Fields(query)
	var tokens []string
	for _, t := range rawTokens {
		tokens = append(tokens, strings.ToLower(t))
	}

	var filtered []Issue
	for _, issue := range issues {
		if matchesAllTokens(issue, tokens) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

func matchesAllTokens(issue Issue, tokens []string) bool {
	issueTitle := strings.ToLower(issue.Title)
	issueID := strings.ToLower(issue.ID)
	issueType := strings.ToLower(string(issue.IssueType))
	issuePriorityLevel := issue.Priority
	issuePriorityName := strings.ToLower(PriorityName(issue.Priority))
	issuePriorityLabel := strings.ToLower(PriorityLabel(issue.Priority))

	for _, token := range tokens {
		matched := false

		switch {
		case strings.HasPrefix(token, "type:"):
			// 1. Check explicit 'type:' tokens
			val := strings.TrimPrefix(token, "type:")
			matched = issueType == val

		case strings.HasPrefix(token, "priority:"):
			// 2. Check explicit 'priority:' tokens
			val := strings.TrimPrefix(token, "priority:")
			switch {
			case val == issuePriorityName:
				matched = true
			case val == "0" && issuePriorityLevel == PriorityCritical:
				matched = true
			case val == "1" && issuePriorityLevel == PriorityHigh:
				matched = true
			case val == "2" && issuePriorityLevel == PriorityMedium:
				matched = true
			case val == "3" && issuePriorityLevel == PriorityLow:
				matched = true
			case val == "4" && issuePriorityLevel == PriorityBacklog:
				matched = true
			}

		case token == issuePriorityLabel:
			// 3. Check shorthand priority (p0, p1, etc)
			matched = true

		default:
			// 4. Free-text search on ID and Title
			matched = strings.Contains(issueID, token) || strings.Contains(issueTitle, token)
		}

		// If any token fails to match, the whole issue fails
		if !matched {
			return false
		}
	}

	return true
}
