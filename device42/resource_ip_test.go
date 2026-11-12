package device42

import "testing"

func TestIsIPAddress(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		wantValid bool
	}{
		{"valid ipv4", "192.168.1.1", true},
		{"valid ipv4 localhost", "127.0.0.1", true},
		{"valid ipv6", "2001:db8::1", true},
		{"valid ipv6 full", "2001:0db8:0000:0000:0000:0000:0000:0001", true},
		{"invalid string", "not-an-ip", false},
		{"invalid ipv4 octet", "192.168.1.999", false},
		{"empty string", "", false},
		{"whitespace", " 192.168.1.1 ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warns, errs := IsIPAddress(tt.value, "ip")
			if len(warns) != 0 {
				t.Fatalf("IsIPAddress(%q) returned unexpected warns: %v", tt.value, warns)
			}
			if got := len(errs) == 0; got != tt.wantValid {
				t.Fatalf("IsIPAddress(%q) valid=%v, want valid=%v (errs=%v)", tt.value, got, tt.wantValid, errs)
			}
		})
	}
}