package gastown

import (
	"fmt"
	"strconv"
	"strings"
)

// Vitals represents parsed output from `gt vitals`.
type Vitals struct {
	Servers []DoltServer
	Backups BackupStatus
	Raw     string // fallback if parsing fails
}

// DoltServer represents a single Dolt server entry from vitals output.
type DoltServer struct {
	Port        string
	Label       string
	PID         int
	DiskUsage   string
	Connections string
	Latency     string
	Running     bool
}

// BackupStatus represents backup freshness from vitals output.
type BackupStatus struct {
	LocalLabel string
	JSONLLabel string
	LocalOK    bool
	JSONLOK    bool
}

// FetchVitals runs `gt vitals` and parses the text output.
func FetchVitals() (*Vitals, error) {
	out, err := runWithTimeout(TimeoutMedium, "gt", "vitals")
	if err != nil {
		return nil, fmt.Errorf("gt vitals: %w", err)
	}
	return ParseVitals(string(out)), nil
}

// ParseVitals parses the text output of `gt vitals` into structured data.
// If parsing fails, the raw text is stored as a fallback.
func ParseVitals(raw string) *Vitals {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return &Vitals{}
	}

	v := &Vitals{Raw: raw}

	// Split into sections by header lines (e.g., "Dolt Servers", "Databases", "Backups")
	sections := splitSections(raw)

	if serverLines, ok := sections["dolt servers"]; ok {
		v.Servers = parseServers(serverLines)
	}
	if backupLines, ok := sections["backups"]; ok {
		v.Backups = parseBackups(backupLines)
	}

	// If we parsed something, clear Raw to signal structured data is available
	if len(v.Servers) > 0 || v.Backups.LocalLabel != "" || v.Backups.JSONLLabel != "" {
		v.Raw = ""
	}

	return v
}

// splitSections splits gt vitals output into named sections.
// Section headers are lines that don't start with whitespace and end with a colon or
// are standalone titles like "Dolt Servers", "Databases", "Backups".
func splitSections(raw string) map[string][]string {
	sections := make(map[string][]string)
	var currentKey string

	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Section headers: non-indented lines (possibly with trailing colon)
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && line[0] != '*' && line[0] != '-' {
			currentKey = strings.ToLower(strings.TrimRight(trimmed, ":"))
			continue
		}

		if currentKey != "" {
			sections[currentKey] = append(sections[currentKey], trimmed)
		}
	}

	return sections
}

// parseServers parses server lines from the "Dolt Servers" section.
// Format varies: bullet may be ●/○ (Unicode) or */- (ASCII), details may be
// on the same line or a continuation line.
func parseServers(lines []string) []DoltServer {
	var servers []DoltServer
	var current *DoltServer

	for _, line := range lines {
		if isBulletLine(line) {
			if current != nil {
				servers = append(servers, *current)
			}
			// ○ (open circle) means stopped; ● (filled) means running
			running := !strings.Contains(line, "○")
			current = &DoltServer{Running: running}
			parseBulletLine(current, line)
			// Also parse detail tokens from the same line (single-line format)
			parseDetailLine(current, line)
			continue
		}

		// Continuation detail line for current server
		if current != nil {
			parseDetailLine(current, line)
		}
	}

	if current != nil {
		servers = append(servers, *current)
	}
	return servers
}

// isBulletLine returns true if the line starts with a bullet character.
func isBulletLine(line string) bool {
	return strings.HasPrefix(line, "*") || strings.HasPrefix(line, "-") ||
		strings.Contains(line, "●") || strings.Contains(line, "○")
}

// parseBulletLine extracts port, label, PID from a server bullet line.
func parseBulletLine(s *DoltServer, line string) {
	// Strip bullet prefix (ASCII and Unicode)
	line = strings.TrimLeft(line, "*- ")
	line = strings.ReplaceAll(line, "●", "")
	line = strings.ReplaceAll(line, "○", "")
	line = strings.TrimSpace(line)
	fields := strings.Fields(line)

	for i, f := range fields {
		if strings.HasPrefix(f, ":") {
			s.Port = f
		} else if strings.EqualFold(f, "PID") && i+1 < len(fields) {
			s.PID, _ = strconv.Atoi(fields[i+1])
		} else if strings.EqualFold(f, "stopped") || strings.EqualFold(f, "down") {
			s.Running = false
		} else if s.Port != "" && s.Label == "" && !strings.EqualFold(f, "PID") && s.PID == 0 {
			// First non-port, non-PID token after port is the label
			if _, err := strconv.Atoi(f); err != nil {
				s.Label = f
			}
		}
	}
}

// parseDetailLine extracts disk usage, connections, latency from a detail line.
func parseDetailLine(s *DoltServer, line string) {
	fields := strings.Fields(line)
	for i, f := range fields {
		lower := strings.ToLower(f)
		if lower == "mb" || lower == "gb" || lower == "kb" {
			if i > 0 {
				s.DiskUsage = fields[i-1] + " " + f
			}
		} else if strings.HasSuffix(lower, "mb") || strings.HasSuffix(lower, "gb") {
			if s.DiskUsage == "" {
				s.DiskUsage = f
			}
		} else if strings.Contains(lower, "conn") {
			if i > 0 {
				s.Connections = fields[i-1] + " " + f
			} else {
				s.Connections = f
			}
		} else if strings.HasSuffix(lower, "ms") || strings.HasSuffix(lower, "µs") {
			s.Latency = f
		}
	}
}

// parseBackups parses the "Backups" section for Local and JSONL status.
func parseBackups(lines []string) BackupStatus {
	var b BackupStatus
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "local:") {
			b.LocalLabel = strings.TrimSpace(line[6:])
			b.LocalOK = isBackupOK(b.LocalLabel)
		} else if strings.HasPrefix(lower, "jsonl:") {
			b.JSONLLabel = strings.TrimSpace(line[6:])
			b.JSONLOK = isBackupOK(b.JSONLLabel)
		}
	}
	return b
}

// isBackupOK returns true if the backup label indicates a healthy state.
func isBackupOK(label string) bool {
	lower := strings.ToLower(label)
	if strings.Contains(lower, "not found") || strings.Contains(lower, "not available") ||
		strings.Contains(lower, "missing") || strings.Contains(lower, "error") ||
		strings.Contains(lower, "stale") || lower == "" {
		return false
	}
	return true
}
