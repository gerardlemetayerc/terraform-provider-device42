package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestPowerCircuitFormDataCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42PowerCircuit().Schema, map[string]interface{}{
		"breakerpanel_id": "30",
		"number":          "1",
		"label":           "tf_label",
		"amps":            20,
		"voltage":         120,
		"connection_name": "tf_conn",
	})

	got := powerCircuitFormData(d, "")

	expected := map[string]string{
		"breakerpanel_id": "30",
		"number":          "1",
		"label":           "tf_label",
		"amps":            "20",
		"voltage":         "120",
		"connection":      "tf_conn",
	}

	if len(got) != len(expected) {
		t.Fatalf("powerCircuitFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("powerCircuitFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}

	if _, ok := got["id"]; ok {
		t.Errorf("powerCircuitFormData() on create must not include an id")
	}
}

func TestPowerCircuitFormDataUpdateIncludesID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42PowerCircuit().Schema, map[string]interface{}{
		"breakerpanel_id": "30",
		"number":          "1",
		"amps":            15,
	})

	got := powerCircuitFormData(d, "87")

	if got["id"] != "87" {
		t.Errorf("powerCircuitFormData() on update must include id, got %q", got["id"])
	}
	if got["amps"] != "15" {
		t.Errorf("powerCircuitFormData() amps = %q, want %q", got["amps"], "15")
	}
	if got["breakerpanel_id"] != "30" {
		t.Errorf("powerCircuitFormData() breakerpanel_id = %q, want %q", got["breakerpanel_id"], "30")
	}
}

func TestPowerCircuitFormDataOmitsEmptyOptionals(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42PowerCircuit().Schema, map[string]interface{}{
		"breakerpanel_id": "30",
		"number":          "1",
		"label":           "",
		"amps":            0,
		"connection_name": "",
	})

	got := powerCircuitFormData(d, "")

	if _, ok := got["label"]; ok {
		t.Errorf("powerCircuitFormData() should omit empty label")
	}
	if _, ok := got["connection"]; ok {
		t.Errorf("powerCircuitFormData() should omit empty connection")
	}
	if got["breakerpanel_id"] != "30" || got["number"] != "1" {
		t.Errorf("powerCircuitFormData() identity fields = %v, want breakerpanel_id=30 number=1", got)
	}
}
