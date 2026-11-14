package device42

import "testing"

// TestProviderInternalValidate validates the whole provider schema (config,
// resources and data sources) using the SDK's internal validator. This catches
// schema definition errors such as invalid field combinations, missing fields
// and inconsistent defaults for every resource and data source at once.
func TestProviderInternalValidate(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider schema is invalid: %s", err)
	}
}
