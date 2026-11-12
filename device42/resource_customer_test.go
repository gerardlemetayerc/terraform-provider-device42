package device42

import "testing"

func TestFilterCustomerByName(t *testing.T) {
	customers := []apiCustomerReadResponse{
		{Name: "Alpha", ID: 1},
		{Name: "Beta", ID: 2},
		{Name: "Alpha", ID: 3},
	}

	tests := []struct {
		name      string
		customers []apiCustomerReadResponse
		search    string
		wantID    int64
		wantErr   bool
	}{
		{"exact match found", customers, "Beta", 2, false},
		{"no match", customers, "Gamma", 0, true},
		{"multiple matches", customers, "Alpha", 0, true},
		{"empty list", []apiCustomerReadResponse{}, "Alpha", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterCustomerByName(tt.customers, tt.search)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("filterCustomerByName(%q) expected error, got none", tt.search)
				}
				return
			}
			if err != nil {
				t.Fatalf("filterCustomerByName(%q) unexpected error: %s", tt.search, err)
			}
			if got.ID != tt.wantID {
				t.Fatalf("filterCustomerByName(%q) ID = %d, want %d", tt.search, got.ID, tt.wantID)
			}
		})
	}
}
