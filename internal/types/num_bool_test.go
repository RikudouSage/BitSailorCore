package types

import (
	"encoding/json"
	"testing"
)

func TestNumBoolJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want NumBool
	}{
		{name: "true", raw: `1`, want: NumBool(true)},
		{name: "false", raw: `0`, want: NumBool(false)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var got NumBool
			if err := json.Unmarshal([]byte(test.raw), &got); err != nil {
				t.Fatalf("Unmarshal() returned error: %v", err)
			}
			if got != test.want {
				t.Fatalf("Unmarshal() = %v, want %v", got, test.want)
			}

			raw, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Marshal() returned error: %v", err)
			}
			if string(raw) != test.raw {
				t.Fatalf("Marshal() = %s, want %s", raw, test.raw)
			}
		})
	}
}

func TestNumBoolUnmarshalRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	var got NumBool
	if err := json.Unmarshal([]byte(`true`), &got); err == nil {
		t.Fatal("Unmarshal() error = nil, want error")
	}
}
