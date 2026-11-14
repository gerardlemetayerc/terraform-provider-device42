package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAutoDiscoveryCloudAccountFormData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryCloudAccount().Schema, map[string]interface{}{
		"name":                         "Openstack Dev",
		"type":                         "OpenStack",
		"server":                       "openstack.example.com",
		"ip":                           "https://192.168.11.80:5000",
		"username":                     "admin",
		"tenant":                       "admin",
		"vendor":                       "openstack",
		"object_category":              "Cloud",
		"regions":                      "region-a,region-b",
		"collect_db_info":              "yes",
		"remove_unfound_instances":     "yes",
		"add_tags_as_custom_fields":    "no",
		"overwrite_object_category":    "no",
		"strip_domain_name":            "no",
		"match_name_only_for_virtuals": "no",
		"schedule_time":                "08:00",
		"schedule_days":                "0,1,2,3,4,5,6",
		"tags":                         "prod,cloud",
	})

	got := autoDiscoveryCloudAccountFormData(d)

	expected := map[string]string{
		"name":                         "Openstack Dev",
		"type":                         "OpenStack",
		"server":                       "openstack.example.com",
		"ip":                           "https://192.168.11.80:5000",
		"username":                     "admin",
		"tenant":                       "admin",
		"vendor":                       "openstack",
		"object_category":              "Cloud",
		"regions":                      "region-a,region-b",
		"collect_db_info":              "yes",
		"remove_unfound_instances":     "yes",
		"add_tags_as_custom_fields":    "no",
		"overwrite_object_category":    "no",
		"strip_domain_name":            "no",
		"match_name_only_for_virtuals": "no",
		"schedule_time":                "08:00",
		"schedule_days":                "0,1,2,3,4,5,6",
		"tags":                         "prod,cloud",
	}

	if len(got) != len(expected) {
		t.Fatalf("autoDiscoveryCloudAccountFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("autoDiscoveryCloudAccountFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
}

func TestAutoDiscoveryCloudAccountFormDataRequiresNameAndType(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42AutoDiscoveryCloudAccount().Schema, map[string]interface{}{
		"name": "Openstack Dev",
		"type": "OpenStack",
	})

	got := autoDiscoveryCloudAccountFormData(d)

	if got["name"] != "Openstack Dev" {
		t.Errorf("autoDiscoveryCloudAccountFormData() name = %q, want Openstack Dev", got["name"])
	}
	if got["type"] != "OpenStack" {
		t.Errorf("autoDiscoveryCloudAccountFormData() type = %q, want OpenStack", got["type"])
	}
	if _, ok := got["server"]; ok {
		t.Errorf("autoDiscoveryCloudAccountFormData() should omit empty server")
	}
	if _, ok := got["regions"]; ok {
		t.Errorf("autoDiscoveryCloudAccountFormData() should omit empty regions")
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("autoDiscoveryCloudAccountFormData() should omit empty tags")
	}
}
