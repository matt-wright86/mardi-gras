package ui

import "testing"

func TestSuperscriptSingleDigit(t *testing.T) {
	cases := map[int]string{
		0: "⁰", 1: "¹", 5: "⁵", 9: "⁹",
	}
	for n, want := range cases {
		if got := Superscript(n); got != want {
			t.Errorf("Superscript(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestSuperscriptMultiDigit(t *testing.T) {
	if got := Superscript(42); got != "⁴²" {
		t.Errorf("Superscript(42) = %q, want ⁴²", got)
	}
	if got := Superscript(100); got != "¹⁰⁰" {
		t.Errorf("Superscript(100) = %q, want ¹⁰⁰", got)
	}
}

func TestSuperscriptNegative(t *testing.T) {
	if got := Superscript(-5); got != "⁰" {
		t.Errorf("Superscript(-5) = %q, want ⁰", got)
	}
}
