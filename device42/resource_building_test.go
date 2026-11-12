package device42

import (
	"encoding/json"
	"testing"
)

func TestMsgID(t *testing.T) {
	tests := []struct {
		name    string
		msg     []interface{}
		want    int
		wantErr bool
	}{
		{"numeric float64 at index 1", []interface{}{"ok", 114.0, "building"}, 114, false},
		{"json.Number at index 1", []interface{}{"ok", json.Number("114"), "building"}, 114, false},
		{"integer at index 1", []interface{}{"ok", 114, "building"}, 114, false},
		{"msg too short", []interface{}{"ok"}, 0, true},
		{"empty msg", nil, 0, true},
		{"non-numeric element", []interface{}{"ok", "114", "building"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgID(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("msgID(%v) err=%v, wantErr=%v", tt.msg, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("msgID(%v) = %d, want %d", tt.msg, got, tt.want)
			}
		})
	}
}

func TestAPIDeleteResponseUnmarshal(t *testing.T) {
	body := `{"deleted":"true","id":"114"}`
	var r apiDeleteResponse
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatalf("unmarshal delete response: %s", err)
	}
	if r.Deleted != "true" {
		t.Fatalf("deleted=%q, want \"true\"", r.Deleted)
	}
	if r.ID != "114" {
		t.Fatalf("id=%q, want \"114\"", r.ID)
	}
}

func TestAPIStringOrNumberUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"json number", `48.85`, "48.85"},
		{"json integer float", `22.0`, "22.0"},
		{"json string", `"60"`, "60"},
		{"json null", `null`, ""},
		{"empty string", `""`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v apiStringOrNumber
			if err := json.Unmarshal([]byte(tt.json), &v); err != nil {
				t.Fatalf("unmarshal %s: %s", tt.json, err)
			}
			if string(v) != tt.want {
				t.Fatalf("unmarshal %s = %q, want %q", tt.json, string(v), tt.want)
			}
		})
	}
}

func TestAPIStringOrNumberUnmarshalInvalid(t *testing.T) {
	var v apiStringOrNumber
	if err := json.Unmarshal([]byte(`{"nested":true}`), &v); err == nil {
		t.Fatal("expected error unmarshaling object into apiStringOrNumber")
	}
}
