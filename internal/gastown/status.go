package gastown

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// TownStatus is the parsed output of `gt status --json`.
type TownStatus struct {
	Agents  []AgentRuntime `json:"agents"`
	Rigs    []RigStatus    `json:"rigs"`
	Convoys []ConvoyInfo   `json:"convoys"`
}

// AgentRuntime represents a single Gas Town agent.
type AgentRuntime struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Rig       string `json:"rig"`
	Running   bool   `json:"running"`
	HasWork   bool   `json:"has_work"`
	WorkTitle string `json:"work_title"`
	HookBead  string `json:"hook_bead"`
	State     string `json:"state"`
	Mail      int    `json:"unread_mail"`
}

// RigStatus represents a Gas Town rig (project).
type RigStatus struct {
	Name   string `json:"name"`
	Agents int    `json:"agents"`
}

// ConvoyInfo represents a Gas Town convoy (delivery batch).
type ConvoyInfo struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Done   int    `json:"done"`
	Total  int    `json:"total"`
}

// FetchStatus runs `gt status --json` and parses the output.
// Returns nil TownStatus (not error) if gt is not available.
func FetchStatus() (*TownStatus, error) {
	out, err := exec.Command("gt", "status", "--json").Output()
	if err != nil {
		return nil, fmt.Errorf("gt status: %w", err)
	}
	var status TownStatus
	if err := json.Unmarshal(out, &status); err != nil {
		return nil, fmt.Errorf("gt status parse: %w", err)
	}
	return &status, nil
}

// AgentForIssue returns the agent working on a given issue, if any.
func (s *TownStatus) AgentForIssue(issueID string) *AgentRuntime {
	if s == nil {
		return nil
	}
	for i := range s.Agents {
		if s.Agents[i].HookBead == issueID {
			return &s.Agents[i]
		}
	}
	return nil
}

// ActiveAgentMap returns issueID -> agent name for all agents with hooked work.
// This bridges Gas Town status to the existing ActiveAgents map[string]string.
func (s *TownStatus) ActiveAgentMap() map[string]string {
	m := make(map[string]string)
	if s == nil {
		return m
	}
	for _, a := range s.Agents {
		if a.HookBead != "" && a.Running {
			m[a.HookBead] = a.Name
		}
	}
	return m
}

// WorkingCount returns the number of agents currently doing work.
func (s *TownStatus) WorkingCount() int {
	n := 0
	if s == nil {
		return n
	}
	for _, a := range s.Agents {
		if a.HasWork && a.Running {
			n++
		}
	}
	return n
}

// UnreadMail returns total unread mail across all agents.
func (s *TownStatus) UnreadMail() int {
	n := 0
	if s == nil {
		return n
	}
	for _, a := range s.Agents {
		n += a.Mail
	}
	return n
}
