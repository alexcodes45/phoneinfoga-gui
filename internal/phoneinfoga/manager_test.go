package phoneinfoga

import "testing"

func TestParseMode(t *testing.T) {
	tests := map[string]Mode{
		"":        ModeServe,
		"serve":   ModeServe,
		"SERVE":   ModeServe,
		"cli":     ModeCLI,
		"CLI":     ModeCLI,
		"unknown": ModeServe,
	}
	for input, want := range tests {
		if got := ParseMode(input); got != want {
			t.Fatalf("ParseMode(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestModeString(t *testing.T) {
	if ModeServe.String() != "serve" {
		t.Fatalf("ModeServe.String() = %q", ModeServe.String())
	}
	if ModeCLI.String() != "cli" {
		t.Fatalf("ModeCLI.String() = %q", ModeCLI.String())
	}
}
