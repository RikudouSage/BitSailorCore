package types

import (
	"encoding/json"
	"testing"
)

func TestCSVSliceJSON(t *testing.T) {
	t.Parallel()

	var values CSVSlice
	if err := json.Unmarshal([]byte(`"one,two,three"`), &values); err != nil {
		t.Fatalf("Unmarshal() returned error: %v", err)
	}
	if values.String() != "one,two,three" {
		t.Fatalf("String() = %q, want one,two,three", values.String())
	}

	raw, err := json.Marshal(values)
	if err != nil {
		t.Fatalf("Marshal() returned error: %v", err)
	}
	if string(raw) != `"one,two,three"` {
		t.Fatalf("Marshal() = %s, want \"one,two,three\"", raw)
	}
}

func TestCSVSliceUnmarshalRejectsNonString(t *testing.T) {
	t.Parallel()

	var values CSVSlice
	if err := json.Unmarshal([]byte(`["one"]`), &values); err == nil {
		t.Fatal("Unmarshal() error = nil, want error")
	}
}

func TestCSVSliceUnmarshalNullLeavesExistingValue(t *testing.T) {
	t.Parallel()

	values := CSVSlice{"existing"}
	if err := json.Unmarshal([]byte(`null`), &values); err != nil {
		t.Fatalf("Unmarshal(null) returned error: %v", err)
	}
	if values.String() != "existing" {
		t.Fatalf("String() = %q, want existing", values.String())
	}
}
