package device42

import "testing"

func TestYesNo(t *testing.T) {
	tests := []struct {
		name string
		b    bool
		want string
	}{
		{"true to yes", true, "yes"},
		{"false to no", false, "no"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yesNo(tt.b); got != tt.want {
				t.Fatalf("yesNo(%v) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
