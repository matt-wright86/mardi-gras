package components

import (
	"strings"
	"testing"

	"github.com/matt-wright86/mardi-gras/internal/gastown"
)

func TestHeaderRigCountMultiRig(t *testing.T) {
	h := Header{
		Width:            120,
		GasTownAvailable: true,
		TownStatus: &gastown.TownStatus{
			Rigs: []gastown.RigStatus{
				{Name: "rig_alpha"},
				{Name: "rig_beta"},
				{Name: "rig_gamma"},
			},
		},
	}
	output := h.View()
	if !strings.Contains(output, "3 rigs") {
		t.Fatalf("expected header to contain '3 rigs' for multi-rig, got: %s", output)
	}
}

func TestHeaderRigCountSingleRig(t *testing.T) {
	h := Header{
		Width:            120,
		GasTownAvailable: true,
		TownStatus: &gastown.TownStatus{
			Rigs: []gastown.RigStatus{
				{Name: "rig_alpha"},
			},
		},
	}
	output := h.View()
	if strings.Contains(output, "rigs") {
		t.Fatalf("expected header to NOT show rig count for single rig, got: %s", output)
	}
}
