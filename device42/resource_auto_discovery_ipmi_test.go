package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryIpmiFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryIpmi().Schema, map[string]interface{}{
		"name":                        "D42 Power SNMP",
		"server":                      "192.168.211.121-192.168.211.160",
		"discovery_type":              "redfish",
		"update_model_if_found":       "yes",
		"overwrite_object_categories": "no",
		"overwrite_device_name":       "no",
		"run_as_operator":             "yes",
		"debug_level":                 "Debug Off",
		"hostname_to_use":             1,
		"tags":                        "prod,power",
	})

	got := autoDiscoveryIpmiFormData(d)

	expected := map[string]string{
		"name":                        "D42 Power SNMP",
		"server":                      "192.168.211.121-192.168.211.160",
		"discovery_type":              "redfish",
		"update_model_if_found":       "yes",
		"overwrite_object_categories": "no",
		"overwrite_device_name":       "no",
		"run_as_operator":             "yes",
		"debug_level":                 "Debug Off",
		"hostname_to_use":             "1",
		"tags":                        "prod,power",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryIpmiFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryIpmiFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryIpmiFormDataRequiresNameAndServer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryIpmi().Schema, map[string]interface{}{
		"name":   "D42 Power SNMP",
		"server": "192.168.211.121-192.168.211.160",
	})

	got := autoDiscoveryIpmiFormData(d)

	if got["name"] != "D42 Power SNMP" {
		t.Errorf("autoDiscoveryIpmiFormData() name = %q, want D42 Power SNMP", got["name"])
	}
	if got["server"] != "192.168.211.121-192.168.211.160" {
		t.Errorf("autoDiscoveryIpmiFormData() server = %q, want 192.168.211.121-192.168.211.160", got["server"])
	}
	if _, ok := got["schedule_time"]; ok {
		t.Errorf("autoDiscoveryIpmiFormData() should omit empty schedule_time")
	}
	if _, ok := got["ip_start"]; ok {
		t.Errorf("autoDiscoveryIpmiFormData() should omit empty ip_start")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryIpmiFormData() should omit empty tags")
	}
}
