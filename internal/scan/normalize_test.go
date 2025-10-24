package scan

import "testing"

func TestNormalizeE164(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := NormalizeE164("(202) 555-0187", "US")
		if err != nil {
			t.Fatalf("NormalizeE164 returned error: %v", err)
		}
		want := "+12025550187"
		if got != want {
			t.Fatalf("want %s, got %s", want, got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if _, err := NormalizeE164("123", "US"); err == nil {
			t.Fatal("expected error for invalid number")
		}
	})

	t.Run("empty", func(t *testing.T) {
		if _, err := NormalizeE164("   ", "US"); err == nil {
			t.Fatal("expected error for empty input")
		}
	})
}
