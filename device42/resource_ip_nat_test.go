package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestIpNatFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42IpNat().Schema, map[string]interface{}{
		"name":             "nat-rule",
		"ip_address_from":  "10.0.0.1",
		"ip_address_to":    "172.16.0.1",
		"notes":            "test nat",
		"protocol":         "tcp",
		"two_way_relation": "True",
	})

	got := ipNatFormData(d)

	expected := map[string]string{
		"name":             "nat-rule",
		"ip_address_from":  "10.0.0.1",
		"ip_address_to":    "172.16.0.1",
		"notes":            "test nat",
		"protocol":         "tcp",
		"two_way_relation": "True",
	}

	if len(got) != len(expected) {
		t.Fatalf("ipNatFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("ipNatFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestIpNatFormDataRequiresFromAndTo(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42IpNat().Schema, map[string]interface{}{
		"name":            "nat-rule",
		"ip_address_from": "10.0.0.1",
		"ip_address_to":   "172.16.0.1",
	})

	got := ipNatFormData(d)

	if got["name"] != "nat-rule" {
		t.Errorf("ipNatFormData() name = %q, want nat-rule", got["name"])
	}
	if got["ip_address_from"] != "10.0.0.1" {
		t.Errorf("ipNatFormData() ip_address_from = %q, want 10.0.0.1", got["ip_address_from"])
	}
	if got["ip_address_to"] != "172.16.0.1" {
		t.Errorf("ipNatFormData() ip_address_to = %q, want 172.16.0.1", got["ip_address_to"])
	}
}
