package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryNetworksFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryNetworks().Schema, map[string]interface{}{
		"name":              "networks_job",
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
		"tags":              "prod,network",
	})

	got := autoDiscoveryNetworksFormData(d)

	expected := map[string]string{
		"name":              "networks_job",
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
		"tags":              "prod,network",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryNetworksFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryNetworksFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryNetworksFormDataRequiresNameAndServer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryNetworks().Schema, map[string]interface{}{
		"name": "networks_job",
	})

	got := autoDiscoveryNetworksFormData(d)

	if got["name"] != "networks_job" {
		t.Errorf("autoDiscoveryNetworksFormData() name = %q, want networks_job", got["name"])
	}
	if _, ok := got["server"]; ok {
		t.Errorf("autoDiscoveryNetworksFormData() should omit empty server")
	}
	if _, ok := got["snmp_retries"]; ok {
		t.Errorf("autoDiscoveryNetworksFormData() should omit empty snmp_retries")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryNetworksFormData() should omit empty tags")
	}
}
