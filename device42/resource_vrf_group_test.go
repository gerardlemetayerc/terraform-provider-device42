package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSuppressTagsAppendDiff(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{"identical", "vrf,test", "vrf,test", true},
		{"same set different order", "test,vrf", "vrf,test", true},
		{"config subset of state (api appended)", "test,updated,vrf,prod", "vrf,updated,prod", true},
		{"config subset of state reversed", "prod,vrf,updated,test", "vrf,updated,prod", true},
		{"missing config tag -> drift", "vrf,test", "vrf,test,extra", false},
		{"empty both", "", "", true},
		{"empty new, non-empty old", "vrf,test", "", false},
		{"empty old, non-empty new", "", "vrf,test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := suppressTagsAppendDiff("tags", tt.old, tt.new, &schema.ResourceData{})
			if got != tt.want {
				t.Errorf("suppressTagsAppendDiff(%q, %q) = %v, want %v", tt.old, tt.new, got, tt.want)
			}
		})
	}
}
