package data

import (
	"testing"
)

func TestSourceLabelJSONL(t *testing.T) {
	tests := []struct {
		name string
		src  Source
		want string
	}{
		{
			name: "JSONL with path",
			src:  Source{Mode: SourceJSONL, Path: "/foo/.beads/issues.jsonl"},
			want: "issues.jsonl",
		},
		{
			name: "JSONL empty path",
			src:  Source{Mode: SourceJSONL},
			want: "issues.jsonl",
		},
		{
			name: "CLI mode",
			src:  Source{Mode: SourceCLI},
			want: "bd list",
		},
		{
			name: "CLI mode ignores path",
			src:  Source{Mode: SourceCLI, Path: "/foo/bar"},
			want: "bd list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.src.Label()
			if got != tt.want {
				t.Errorf("Source.Label() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCheckBdVersionKnownBroken(t *testing.T) {
	got := parseBdVersionWarning("bd version 0.59.0")
	if got == "" {
		t.Fatal("expected warning for v0.59.0, got empty string")
	}
	if got != "bd v0.59.0 has a known bug where --json is ignored; upgrade to v0.59.1+" {
		t.Errorf("unexpected warning: %q", got)
	}
}

func TestCheckBdVersionOK(t *testing.T) {
	got := parseBdVersionWarning("bd version 0.58.0")
	if got != "" {
		t.Errorf("expected no warning for v0.58.0, got %q", got)
	}
}

func TestCheckBdVersionUnparseable(t *testing.T) {
	cases := []string{
		"",
		"garbled output here",
		"bd",
		"\x00\xff",
	}
	for _, input := range cases {
		got := parseBdVersionWarning(input)
		if got != "" {
			t.Errorf("parseBdVersionWarning(%q) = %q, want empty", input, got)
		}
	}
}
