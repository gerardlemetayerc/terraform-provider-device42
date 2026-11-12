package device42

import (
	"encoding/json"
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

func TestAPICodeUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    apiCode
		wantErr bool
	}{
		{"json number", `0`, 0, false},
		{"json number non-zero", `400`, 400, false},
		{"quoted string", `"0"`, 0, false},
		{"quoted string non-zero", `"400"`, 400, false},
		{"invalid string", `"abc"`, 0, true},
		{"invalid type", `true`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c apiCode
			err := json.Unmarshal([]byte(tt.json), &c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json %s: err=%v, wantErr=%v", tt.json, err, tt.wantErr)
			}
			if !tt.wantErr && c != tt.want {
				t.Fatalf("json %s: code=%d, want=%d", tt.json, c, tt.want)
			}
		})
	}
}

func TestArchiveResponseDecodesStringCode(t *testing.T) {
	body := `{"msg":["Object archived",46,"db-080-westport",true,true],"code":"0"}`
	var r apiArchiveResponse
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatalf("unmarshal archive response: %s", err)
	}
	if r.Code != 0 {
		t.Fatalf("code=%d, want 0", r.Code)
	}
	if len(r.Msg) != 5 {
		t.Fatalf("msg has %d elements, want 5", len(r.Msg))
	}
}