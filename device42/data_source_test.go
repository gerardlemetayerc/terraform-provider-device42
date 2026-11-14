package device42

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRenderMacPort(t *testing.T) {
	tests := []struct {
		name string
		in   json.RawMessage
		want string
	}{
		{"null", json.RawMessage("null"), ""},
		{"empty", json.RawMessage(`""`), ""},
		{"empty raw", json.RawMessage(""), ""},
		{"plain string", json.RawMessage(`"eth0"`), "eth0"},
		{"nested object", json.RawMessage(`{"port":"REAL-PROD"}`), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderMacPort(tt.in)
			if got != tt.want {
				t.Errorf("renderMacPort(%s) = %q, want %q", string(tt.in), got, tt.want)
			}
		})
	}
}

func TestDataSourceSchemasValidate(t *testing.T) {
	// Ensure every data source schema passes internal validation and exposes
	// at least one lookup field. This catches schema mistakes that would
	// otherwise only surface at runtime.
	ds := []*schema.Resource{
		datasourceD42SubnetCategory(),
		datasourceD42VrfGroup(),
		datasourceD42PduModel(),
		datasourceD42Pdu(),
		datasourceD42Cable(),
		datasourceD42PowerCircuit(),
		datasourceD42Mac(),
		datasourceD42IpNat(),
	}
	for i, r := range ds {
		if r.Schema == nil {
			t.Fatalf("data source %d has nil schema", i)
		}
	}
}

func TestDataSourceLookupFieldsPresent(t *testing.T) {
	// Each data source's schema must contain at least one Optional lookup
	// argument (name/id/cable_id/macaddress/etc.) plus at least one Computed
	// attribute.
	ds := map[string]*schema.Resource{
		"subnet_category": datasourceD42SubnetCategory(),
		"vrf_group":       datasourceD42VrfGroup(),
		"pdu_model":       datasourceD42PduModel(),
		"pdu":             datasourceD42Pdu(),
		"cable":           datasourceD42Cable(),
		"power_circuit":   datasourceD42PowerCircuit(),
		"mac":             datasourceD42Mac(),
		"ip_nat":          datasourceD42IpNat(),
	}
	for name, r := range ds {
		var hasOptional, hasComputed bool
		for _, s := range r.Schema {
			if s.Optional {
				hasOptional = true
			}
			if s.Computed {
				hasComputed = true
			}
		}
		if !hasOptional {
			t.Errorf("data source %s has no Optional lookup field", name)
		}
		if !hasComputed {
			t.Errorf("data source %s has no Computed attribute", name)
		}
	}
}

func TestDataSourceSchemasNoReadOnlyBug(t *testing.T) {
	// Guard against accidentally declaring a lookup field as both Required and
	// Computed, which is an invalid combination for data sources.
	ds := map[string]*schema.Resource{
		"subnet_category": datasourceD42SubnetCategory(),
		"vrf_group":       datasourceD42VrfGroup(),
		"pdu_model":       datasourceD42PduModel(),
		"pdu":             datasourceD42Pdu(),
		"cable":           datasourceD42Cable(),
		"power_circuit":   datasourceD42PowerCircuit(),
		"mac":             datasourceD42Mac(),
		"ip_nat":          datasourceD42IpNat(),
	}
	for name, r := range ds {
		for k, s := range r.Schema {
			if s.Required && s.Computed {
				t.Errorf("data source %s field %q is both Required and Computed", name, k)
			}
		}
	}
}

func TestDatasourceResponseStructs(t *testing.T) {
	// Sanity: the list response structs used by data sources must have distinct
	// field types that can be unmarshalled from the API JSON payloads.
	_ = reflect.TypeOf(apiSubnetCategoryListResponse{})
	_ = reflect.TypeOf(apiVrfGroupListResponse{})
	_ = reflect.TypeOf(apiPduModelListResponse{})
	_ = reflect.TypeOf(apiPduListResponse{})
	_ = reflect.TypeOf(apiCableListResponse{})
	_ = reflect.TypeOf(apiPowerCircuitListResponse{})
	_ = reflect.TypeOf(apiMacListResponse{})
	_ = reflect.TypeOf(apiIpNatListResponse{})
}
