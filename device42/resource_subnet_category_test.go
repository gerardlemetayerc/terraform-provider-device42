package device42

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSubnetCategoryFormDataCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42SubnetCategory().Schema, map[string]interface{}{
		"name":        "myCategory",
		"description": "my subnet category",
		"groups":      "prod,data",
	})

	got := subnetCategoryFormData(d, "")

	expected := map[string]string{
		"name":        "myCategory",
		"description": "my subnet category",
		"groups":      "prod,data",
	}

	if len(got) != len(expected) {
		t.Fatalf("subnetCategoryFormData() returned %d fields, want %d: %v", len(got), len(expected), got)
	}
	for k, want := range expected {
		if got[k] != want {
			t.Errorf("subnetCategoryFormData()[%q] = %q, want %q", k, got[k], want)
		}
	}
	if _, ok := got["id"]; ok {
		t.Errorf("subnetCategoryFormData() on create must not include an id")
	}
}

func TestSubnetCategoryFormDataUpdateIncludesID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42SubnetCategory().Schema, map[string]interface{}{
		"name": "myCategory",
	})

	got := subnetCategoryFormData(d, "12")

	if got["id"] != "12" {
		t.Errorf("subnetCategoryFormData() on update must include id, got %q", got["id"])
	}
	if got["name"] != "myCategory" {
		t.Errorf("subnetCategoryFormData() name = %q, want myCategory", got["name"])
	}
}

func TestSubnetCategoryFormDataOmitsEmptyOptionals(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceD42SubnetCategory().Schema, map[string]interface{}{
		"name":        "myCategory",
		"description": "",
		"groups":      "",
	})

	got := subnetCategoryFormData(d, "")

	if _, ok := got["description"]; ok {
		t.Errorf("subnetCategoryFormData() should omit empty description")
	}
	if _, ok := got["groups"]; ok {
		t.Errorf("subnetCategoryFormData() should omit empty groups")
	}
	if got["name"] != "myCategory" {
		t.Errorf("subnetCategoryFormData() name = %q, want myCategory", got["name"])
	}
}
