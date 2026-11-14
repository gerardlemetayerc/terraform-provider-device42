package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryPingSweepFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryPingSweep().Schema, map[string]interface{}{
		"name":                        "Adhoc Job created on 2015-11-19 16:07:01.712790",
		"networks":                    "8.8.8.0/28\n10.0.0.0/24",
		"create_new_subnet":           "yes",
		"debug_level":                 "Debug Off",
		"overwrite_subnet_categories": "no",
		"reverse_dns":                 "no",
		"strip_domain_suffix":         "no",
		"schedule_time":               "",
	})

	got := autoDiscoveryPingSweepFormData(d)

	expected := map[string]string{
		"name":                        "Adhoc Job created on 2015-11-19 16:07:01.712790",
		"networks":                    "8.8.8.0/28\n10.0.0.0/24",
		"create_new_subnet":           "yes",
		"debug_level":                 "Debug Off",
		"overwrite_subnet_categories": "no",
		"reverse_dns":                 "no",
		"strip_domain_suffix":         "no",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryPingSweepFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryPingSweepFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryPingSweepFormDataRequiresNameAndNetworks(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryPingSweep().Schema, map[string]interface{}{
		"name":     "Adhoc Job",
		"networks": "8.8.8.0/28",
	})

	got := autoDiscoveryPingSweepFormData(d)

	if got["name"] != "Adhoc Job" {
		t.Errorf("autoDiscoveryPingSweepFormData() name = %q, want Adhoc Job", got["name"])
	}
	if got["networks"] != "8.8.8.0/28" {
		t.Errorf("autoDiscoveryPingSweepFormData() networks = %q, want 8.8.8.0/28", got["networks"])
	}
	if _, ok := got["schedule_time"]; ok {
		t.Errorf("autoDiscoveryPingSweepFormData() should omit empty schedule_time")
	}
	if _, ok := got["subnet_category"]; ok {
		t.Errorf("autoDiscoveryPingSweepFormData() should omit empty subnet_category")
	}
	if _, ok := got["debug_level"]; ok {
		t.Errorf("autoDiscoveryPingSweepFormData() should omit empty debug_level")
	}
}
