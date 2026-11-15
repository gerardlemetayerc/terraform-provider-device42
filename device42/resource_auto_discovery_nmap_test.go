package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryNmapFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryNmap().Schema, map[string]interface{}{
		"name":            "nmap_scan",
		"server":          "scanner.domain.local",
		"start_ip":        "10.0.0.1",
		"end_ip":          "10.0.0.100",
		"target":          "10.0.0.0/24",
		"exclude_servers": "10.99.0.2",
		"schedule_time":   "02:00",
		"schedule_days":   "0,1,2,3,4,5,6",
		"tags":            "prod,scan",
	})

	got := autoDiscoveryNmapFormData(d)

	expected := map[string]string{
		"name":            "nmap_scan",
		"server":          "scanner.domain.local",
		"start_ip":        "10.0.0.1",
		"end_ip":          "10.0.0.100",
		"target":          "10.0.0.0/24",
		"exclude_servers": "10.99.0.2",
		"schedule_time":   "02:00",
		"schedule_days":   "0,1,2,3,4,5,6",
		"tags":            "prod,scan",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryNmapFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryNmapFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryNmapFormDataRequiresNameAndServer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryNmap().Schema, map[string]interface{}{
		"name": "nmap_scan",
	})

	got := autoDiscoveryNmapFormData(d)

	if got["name"] != "nmap_scan" {
		t.Errorf("autoDiscoveryNmapFormData() name = %q, want nmap_scan", got["name"])
	}
	if _, ok := got["server"]; ok {
		t.Errorf("autoDiscoveryNmapFormData() should omit empty server")
	}
	if _, ok := got["target"]; ok {
		t.Errorf("autoDiscoveryNmapFormData() should omit empty target")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryNmapFormData() should omit empty tags")
	}
}
