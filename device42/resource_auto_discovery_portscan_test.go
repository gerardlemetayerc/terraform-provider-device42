package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryPortScanFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryPortScan().Schema, map[string]interface{}{
		"name":                        "tcp_scan",
		"server":                      "10.99.0.1",
		"start_ip":                    "10.99.0.1",
		"end_ip":                      "10.99.0.254",
		"capture_hosts_file":          "yes",
		"discover_cluster":            "no",
		"overwrite_device_categories": "yes",
		"schedule_time":               "08:00",
		"schedule_days":               "0,1,2,3,4,5,6",
		"tags":                        "prod,scan",
	})

	got := autoDiscoveryPortScanFormData(d)

	expected := map[string]string{
		"name":                        "tcp_scan",
		"server":                      "10.99.0.1",
		"start_ip":                    "10.99.0.1",
		"end_ip":                      "10.99.0.254",
		"capture_hosts_file":          "yes",
		"discover_cluster":            "no",
		"overwrite_device_categories": "yes",
		"schedule_time":               "08:00",
		"schedule_days":               "0,1,2,3,4,5,6",
		"tags":                        "prod,scan",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryPortScanFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryPortScanFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryPortScanFormDataRequiresName(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryPortScan().Schema, map[string]interface{}{
		"name":   "tcp_scan",
		"server": "10.99.0.1",
	})

	got := autoDiscoveryPortScanFormData(d)

	if got["name"] != "tcp_scan" {
		t.Errorf("autoDiscoveryPortScanFormData() name = %q, want tcp_scan", got["name"])
	}
	if got["server"] != "10.99.0.1" {
		t.Errorf("autoDiscoveryPortScanFormData() server = %q, want 10.99.0.1", got["server"])
	}
	if _, ok := got["start_ip"]; ok {
		t.Errorf("autoDiscoveryPortScanFormData() should omit empty start_ip")
	}
	if _, ok := got["capture_hosts_file"]; ok {
		t.Errorf("autoDiscoveryPortScanFormData() should omit empty capture_hosts_file")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryPortScanFormData() should omit empty tags")
	}
}
