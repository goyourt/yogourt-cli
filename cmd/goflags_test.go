package cmd

import "testing"

func TestWithStripDwarfLinkerFlag(t *testing.T) {
	cases := map[string]string{
		"":                         "-ldflags=-w",
		"-mod=mod":                 "-mod=mod -ldflags=-w",
		"-ldflags=-s":              "-ldflags=-s",
		"-mod=mod -ldflags=-X a=b": "-mod=mod -ldflags=-X a=b",
	}

	for goflags, want := range cases {
		if got := withStripDwarfLinkerFlag(goflags); got != want {
			t.Errorf("withStripDwarfLinkerFlag(%q) = %q, want %q", goflags, got, want)
		}
	}
}
