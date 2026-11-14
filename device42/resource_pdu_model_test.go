package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestPduModelFormDataCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42PduModel().Schema, map[string]interface{}{
		"name":         "AP7800",
		"manufacturer": "APC",
		"description":  "a pdu model",
		"notes":        "note",
		"depth":        "Half Depth",
		"size":         42,
		"width":        48,
		"height":       44,
	})

	got := pduModelFormData(d, "")

	expected := map[string]string{
		"name":         "AP7800",
		"manufacturer": "APC",
		"description":  "a pdu model",
		"notes":        "note",
		"depth":        "Half Depth",
		"size":         "42",
		"width":        "48",
		"height":       "44",
	}

	if len(got) != len(expected) {
		t.Fatalf("pduModelFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("pduModelFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
	if _, ok := got["id"]; ok {
		t.Errorf("pduModelFormData() on create must not include an id")
	}
}

func TestPduModelFormDataUpdateIncludesID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42PduModel().Schema, map[string]interface{}{
		"name":         "AP7800",
		"manufacturer": "APC",
	})

	got := pduModelFormData(d, "19")

	if got["id"] != "19" {
		t.Errorf("pduModelFormData() on update must include id, got %q", got["id"])
	}
	if got["manufacturer"] != "APC" {
		t.Errorf("pduModelFormData() manufacturer = %q, want APC", got["manufacturer"])
	}
}

func TestPduModelFormDataOmitsEmptyOptionals(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42PduModel().Schema, map[string]interface{}{
		"name":         "AP7800",
		"manufacturer": "",
		"notes":        "",
	})

	got := pduModelFormData(d, "")

	if _, ok := got["manufacturer"]; ok {
		t.Errorf("pduModelFormData() should omit empty manufacturer")
	}
	if _, ok := got["notes"]; ok {
		t.Errorf("pduModelFormData() should omit empty notes")
	}
	if got["name"] != "AP7800" {
		t.Errorf("pduModelFormData() name = %q, want AP7800", got["name"])
	}
}
