package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryDnsFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryDns().Schema, map[string]interface{}{
		"zonename":      "rajlog.pvt",
		"nameserver":    "192.168.11.104",
		"schedule_time": "10:32",
		"schedule_days": "0,1,2,3,4,5,6",
	})

	got := autoDiscoveryDnsFormData(d)

	expected := map[string]string{
		"zonename":      "rajlog.pvt",
		"nameserver":    "192.168.11.104",
		"schedule_time": "10:32",
		"schedule_days": "0,1,2,3,4,5,6",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryDnsFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryDnsFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryDnsFormDataRequiresZoneNameAndNameserver(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryDns().Schema, map[string]interface{}{
		"zonename":   "rajlog.pvt",
		"nameserver": "192.168.11.104",
	})

	got := autoDiscoveryDnsFormData(d)

	if got["zonename"] != "rajlog.pvt" {
		t.Errorf("autoDiscoveryDnsFormData() zonename = %q, want rajlog.pvt", got["zonename"])
	}
	if got["nameserver"] != "192.168.11.104" {
		t.Errorf("autoDiscoveryDnsFormData() nameserver = %q, want 192.168.11.104", got["nameserver"])
	}
	if _, ok := got["schedule_time"]; ok {
		t.Errorf("autoDiscoveryDnsFormData() should omit empty schedule_time")
	}
	if _, ok := got["schedule_days"]; ok {
		t.Errorf("autoDiscoveryDnsFormData() should omit empty schedule_days")
	}
}
