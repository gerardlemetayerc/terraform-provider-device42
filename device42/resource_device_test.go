package device42

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFlattenCustomFields(t *testing.T) {
	tests := []struct {
		name string
		in   []customField
		want map[string]interface{}
	}{
		{
			name: "empty input",
			in:   nil,
			want: map[string]interface{}{},
		},
		{
			name: "single field",
			in:   []customField{{Key: "owner", Value: "team-a"}},
			want: map[string]interface{}{"owner": "team-a"},
		},
		{
			name: "multiple mixed value types",
			in: []customField{
				{Key: "owner", Value: "team-a"},
				{Key: "priority", Value: 1},
				{Key: "enabled", Value: true},
			},
			want: map[string]interface{}{
				"owner":    "team-a",
				"priority": 1,
				"enabled":  true,
			},
		},
		{
			name: "duplicate keys keep last value",
			in: []customField{
				{Key: "owner", Value: "team-a"},
				{Key: "owner", Value: "team-b"},
			},
			want: map[string]interface{}{"owner": "team-b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenCustomFields(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("flattenCustomFields(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestSuppressCustomFieldsDiffs(t *testing.T) {
	resSchema := map[string]*schema.Schema{
		"custom_fields": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}

	t.Run("field present in state is not suppressed", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resSchema, map[string]interface{}{
			"custom_fields": map[string]interface{}{"foo": "bar"},
		})
		if suppressCustomFieldsDiffs("custom_fields.foo", "old", "new", d) {
			t.Fatal("expected no suppression for a field already in state")
		}
	})

	t.Run("field absent from state is suppressed", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resSchema, map[string]interface{}{
			"custom_fields": map[string]interface{}{"foo": "bar"},
		})
		if !suppressCustomFieldsDiffs("custom_fields.baz", "old", "new", d) {
			t.Fatal("expected suppression for a field not in state")
		}
	})
}