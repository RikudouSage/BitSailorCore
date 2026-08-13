package helper

import "testing"

func TestSliceToLookupMap(t *testing.T) {
	t.Parallel()

	got := SliceToLookupMap([]string{"one", "two", "one"})
	if len(got) != 2 {
		t.Fatalf("len(map) = %d, want 2", len(got))
	}
	if !got["one"] || !got["two"] {
		t.Fatalf("map = %#v, want one and two present", got)
	}
}

func TestSliceToLookupMapNil(t *testing.T) {
	t.Parallel()

	got := SliceToLookupMap[string](nil)
	if got == nil {
		t.Fatal("SliceToLookupMap(nil) = nil, want empty map")
	}
	if len(got) != 0 {
		t.Fatalf("len(map) = %d, want 0", len(got))
	}
}
