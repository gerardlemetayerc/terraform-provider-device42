package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMacFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42Mac().Schema, map[string]interface{}{
		"macaddress": "00:11:22:33:44:55",
		"vlan_id":    "14",
		"port":       "eth0",
	})

	got := macFormData(d)

	expected := map[string]string{
		"macaddress": "00:11:22:33:44:55",
		"vlan_id":    "14",
		"port":       "eth0",
	}

	if len(got) != len(expected) {
		t.Fatalf("macFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("macFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestMacFormDataOmitsEmptyOptionals(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42Mac().Schema, map[string]interface{}{
		"macaddress": "00:11:22:33:44:55",
		"vlan_id":    "",
		"port":       "",
	})

	got := macFormData(d)

	if _, ok := got["vlan_id"]; ok {
		t.Errorf("macFormData() should omit empty vlan_id")
	}
	if _, ok := got["port"]; ok {
		t.Errorf("macFormData() should omit empty port")
	}
	if got["macaddress"] != "00:11:22:33:44:55" {
		t.Errorf("macFormData() macaddress = %q, want 00:11:22:33:44:55", got["macaddress"])
	}
}
