package device42

import "testing"

func TestFindContact(t *testing.T) {
	customer := customerWithContacts{
		apiCustomerReadResponse: apiCustomerReadResponse{Name: "Acme", ID: 1},
		Contacts: []apiCustomerContactResponse{
			{Name: "Alice", Type: "Technical", Email: "alice@example.com", Phone: "+1", ID: 100},
			{Name: "Bob", Type: "Administrative", Address: "Foo St", ID: 101},
		},
	}

	tests := []struct {
		name      string
		contactID int
		wantID    int64
		wantOK    bool
	}{
		{"exact id match", 100, 100, true},
		{"id not found", 999, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := findContact(customer, tt.contactID)
			if ok != tt.wantOK {
				t.Fatalf("findContact(%d) ok = %v, want %v", tt.contactID, ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got.ID != tt.wantID {
				t.Fatalf("findContact(%d) ID = %d, want %d", tt.contactID, got.ID, tt.wantID)
			}
		})
	}
}

func TestFindContactByName(t *testing.T) {
	customer := customerWithContacts{
		apiCustomerReadResponse: apiCustomerReadResponse{Name: "Acme", ID: 1},
		Contacts: []apiCustomerContactResponse{
			{Name: "Alice", Type: "Technical", ID: 100},
			{Name: "Bob", Type: "Administrative", ID: 101},
			{Name: "Alice", Type: "Billing", ID: 102},
		},
	}

	tests := []struct {
		name        string
		contactName string
		wantCount   int
		wantErr     bool
	}{
		{"name match", "Bob", 1, false},
		{"name not found", "Carol", 0, true},
		{"ambiguous name", "Alice", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findContactByName(customer, tt.contactName)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("findContactByName(%q) expected error, got none", tt.contactName)
				}
				return
			}
			if err != nil {
				t.Fatalf("findContactByName(%q) unexpected error: %s", tt.contactName, err)
			}
			if len(got) != tt.wantCount {
				t.Fatalf("findContactByName(%q) got %d matches, want %d", tt.contactName, len(got), tt.wantCount)
			}
		})
	}
}
