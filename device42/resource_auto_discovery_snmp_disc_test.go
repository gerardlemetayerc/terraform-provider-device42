package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoverySnmpDiscFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoverySnmpDisc().Schema, map[string]interface{}{
		"name":              "snmp_job",
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
		"tags":              "prod,snmp",
	})

	got := autoDiscoverySnmpDiscFormData(d)

	expected := map[string]string{
		"name":              "snmp_job",
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
		"tags":              "prod,snmp",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoverySnmpDiscFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoverySnmpDiscFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoverySnmpDiscFormDataRequiresNameAndServer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoverySnmpDisc().Schema, map[string]interface{}{
		"name": "snmp_job",
	})

	got := autoDiscoverySnmpDiscFormData(d)

	if got["name"] != "snmp_job" {
		t.Errorf("autoDiscoverySnmpDiscFormData() name = %q, want snmp_job", got["name"])
	}
	if _, ok := got["server"]; ok {
		t.Errorf("autoDiscoverySnmpDiscFormData() should omit empty server")
	}
	if _, ok := got["snmp_retries"]; ok {
		t.Errorf("autoDiscoverySnmpDiscFormData() should omit empty snmp_retries")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoverySnmpDiscFormData() should omit empty tags")
	}
}
