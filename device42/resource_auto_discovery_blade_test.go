package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryBladeFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryBlade().Schema, map[string]interface{}{
		"name":              "blade_job",
		"server":            "scanner.domain.local",
		"start_ip":          "10.0.0.1",
		"end_ip":            "10.0.0.100",
		"exclude_servers":   "10.99.0.2",
		"snmp_retries":      3,
		"snmp_timeout":      30,
		"snmp_min_parallel": 50,
		"do_port_check":     "true",
		"clear_customer":    "false",
		"schedule_time":     "03:00",
		"schedule_days":     "1,2,3",
		"snmp_v3":           "authpriv",
		"tags":              "prod,blade",
	})

	got := autoDiscoveryBladeFormData(d)

	expected := map[string]string{
		"name":              "blade_job",
		"server":            "scanner.domain.local",
		"start_ip":          "10.0.0.1",
		"end_ip":            "10.0.0.100",
		"exclude_servers":   "10.99.0.2",
		"snmp_retries":      "3",
		"snmp_timeout":      "30",
		"snmp_min_parallel": "50",
		"do_port_check":     "true",
		"clear_customer":    "false",
		"schedule_time":     "03:00",
		"schedule_days":     "1,2,3",
		"snmp_v3":           "authpriv",
		"tags":              "prod,blade",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryBladeFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryBladeFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryBladeFormDataRequiresNameAndServer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryBlade().Schema, map[string]interface{}{
		"name": "blade_job",
	})

	got := autoDiscoveryBladeFormData(d)

	if got["name"] != "blade_job" {
		t.Errorf("autoDiscoveryBladeFormData() name = %q, want blade_job", got["name"])
	}
	if _, ok := got["server"]; ok {
		t.Errorf("autoDiscoveryBladeFormData() should omit empty server")
	}
	if _, ok := got["snmp_retries"]; ok {
		t.Errorf("autoDiscoveryBladeFormData() should omit empty snmp_retries")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryBladeFormData() should omit empty tags")
	}
}
