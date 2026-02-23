package data

import "strings"

// ExternalRef represents a parsed cross-rig external dependency reference.
// Format: "external:<rig>:<id>" — e.g., "external:gastown:gt-c3f2"
type ExternalRef struct {
	Rig      string // source rig name (e.g., "gastown", "wyvern")
	IssueID  string // issue ID within that rig
	Original string // the full original string
}

// ParseExternalRef extracts rig and issue ID from an external reference.
// Returns nil if the string is not an external reference.
func ParseExternalRef(ref string) *ExternalRef {
	if !strings.HasPrefix(ref, "external:") {
		return nil
	}
	parts := strings.SplitN(ref, ":", 3)
	if len(parts) != 3 {
		return nil
	}
	return &ExternalRef{
		Rig:      parts[1],
		IssueID:  parts[2],
		Original: ref,
	}
}

// CrossRigDeps extracts cross-rig dependencies from an issue's dependency list.
func CrossRigDeps(issue *Issue) []ExternalRef {
	var refs []ExternalRef
	for _, dep := range issue.Dependencies {
		if ref := ParseExternalRef(dep.DependsOnID); ref != nil {
			refs = append(refs, *ref)
		}
	}
	return refs
}

// CrossRigSummary computes a summary of cross-rig dependencies across all issues.
// Returns a map of rig name → count of dependencies to/from that rig.
func CrossRigSummary(issues []Issue) map[string]int {
	rigs := make(map[string]int)
	for i := range issues {
		for _, dep := range issues[i].Dependencies {
			if ref := ParseExternalRef(dep.DependsOnID); ref != nil {
				rigs[ref.Rig]++
			}
		}
	}
	return rigs
}
