package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryCertificateFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryCertificate().Schema, map[string]interface{}{
		"name":              "cert_hostname",
		"server":            "hostname.domain.local",
		"ports":             "443,8443",
		"exclude_servers":   "10.99.0.2",
		"follow_chain":      "yes",
		"get_cipher_suites": "yes",
		"schedule_time":     "08:00",
		"schedule_days":     "0,1,2,3,4,5,6",
		"tags":              "prod,ssl",
	})

	got := autoDiscoveryCertificateFormData(d)

	expected := map[string]string{
		"name":              "cert_hostname",
		"server":            "hostname.domain.local",
		"exclude_servers":   "10.99.0.2",
		"ports":             "443,8443",
		"get_cipher_suites": "yes",
		"follow_chain":      "yes",
		"schedule_time":     "08:00",
		"schedule_days":     "0,1,2,3,4,5,6",
		"tags":              "prod,ssl",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryCertificateFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryCertificateFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryCertificateFormDataRequiresNameAndServer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryCertificate().Schema, map[string]interface{}{
		"name":   "cert_hostname",
		"server": "hostname.domain.local",
	})

	got := autoDiscoveryCertificateFormData(d)

	if got["name"] != "cert_hostname" {
		t.Errorf("autoDiscoveryCertificateFormData() name = %q, want cert_hostname", got["name"])
	}
	if got["server"] != "hostname.domain.local" {
		t.Errorf("autoDiscoveryCertificateFormData() server = %q, want hostname.domain.local", got["server"])
	}
	if _, ok := got["ports"]; ok {
		t.Errorf("autoDiscoveryCertificateFormData() should omit empty ports")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryCertificateFormData() should omit empty tags")
	}
}
