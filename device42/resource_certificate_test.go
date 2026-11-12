package device42

import "testing"

func TestFilterCertificateByName(t *testing.T) {
	certs := []apiCertificateDetails{
		{ID: 1, IssuedTo: "*.example.com"},
		{ID: 2, IssuedTo: "*.other.com"},
		{ID: 3, IssuedTo: "*.example.com"},
	}

	tests := []struct {
		name    string
		certs   []apiCertificateDetails
		search  string
		wantID  int64
		wantErr bool
	}{
		{"exact match found", certs, "*.other.com", 2, false},
		{"no match", certs, "*.missing.com", 0, true},
		{"multiple matches", certs, "*.example.com", 0, true},
		{"empty list", []apiCertificateDetails{}, "*.example.com", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterCertificateByName(tt.certs, tt.search)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("filterCertificateByName(%q) expected error, got none", tt.search)
				}
				return
			}
			if err != nil {
				t.Fatalf("filterCertificateByName(%q) unexpected error: %s", tt.search, err)
			}
			if got.ID != tt.wantID {
				t.Fatalf("filterCertificateByName(%q) ID = %d, want %d", tt.search, got.ID, tt.wantID)
			}
		})
	}
}
